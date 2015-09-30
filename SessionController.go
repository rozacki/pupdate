package main

import(
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"strings"
	_"os"
)

const(
	FormatStruct	=	"%#v"
)

type SessionController struct {
	Configuration SessionConfiguration
}

/*
//simple string error type
type StringError struct{
	string
}
*/
func (this StringError) Error() string{
	return string(this)
}

type StringError string

func makeSessionController(configuration SessionConfiguration)(session* SessionController,err error){
	//is there  any task to schedule?
	if len(configuration.Tasks)==0{
		return nil,StringError("no tasks to schedule, stop")
	}

	//generate session name
	sessionController:=SessionController{Configuration:configuration}
	//open logging file explicitly
	if err:= GLogging.OpenLog();err!=nil{
		return nil,err
	}

	return &sessionController,nil
}

//Runs job in order
func (this *SessionController) StartTasks() {
	for _,task:=range this.Configuration.Tasks{
		//if task is disabled then we skip it
		if task.Disabled{
			this.Configuration.TaskDisabled(task.Name)
			continue
		}
		//
		this.Configuration.TaskCounter++

		//synchronous task
		taskData:=this.startTask(task)

		//on first task failed terminates session
		if taskData.Failed{
			this.Configuration.SessionFailed(fmt.Sprintf("'%s'. Task details %+v",taskData.LastDroppedJob.LastErrorMsg,taskData))
			return
		}

		//should I finish
		select{
		//this channel is used as global flag to finish processing. it is used by main thread and allow user to signal termination
		case <-GDone:
			fmt.Printf("session terminated\n")
			return
		case  <-time.After(1* time.Millisecond):
		}
	}
	this.Configuration.SessionSuccess()
}
//runs single task by spawning jobs
//todo: cancellation via done channel
//todo: move to taskcontroller
func (this *SessionController) startTask(configuration TaskConfiguration) TaskData {
	configuration.EventStartTask()

	//create new task data structure where we hold op data and configuration
	var taskData 			= 	TaskData{Id: this.Configuration.TaskCounter, DataCursor: configuration.Min, Status: "working", TaskConfiguration: configuration,CreationTime:time.Now()}
	var JobsCount	uint64
	if configuration.Step<=0 || configuration.Max<=configuration.Min || configuration.Concurrency<1{
		JobsCount					=1
		configuration.Step			=0
		configuration.Max			=0
		configuration.Min			=0
		configuration.Concurrency	=1
	}else {
		JobsCount	= ((configuration.Max - configuration.Min)/configuration.Step)
	}
	//when job finished then posts its status here
	var jobDataChannel		=	make(chan *JobData, configuration.Concurrency)
	//
	defer func(){
		if taskData.Failed {
			//record that fact
			configuration.EventFailTask(taskData.LastDroppedJob.LastErrorMsg)
		}else{
			configuration.RowsAffected(taskData.RowsAffected)
			configuration.EventSuccessTask()
		}
		//last serialisation
		SerialiseStruct(this.Configuration)
		//close task<-job comm channel
		close(jobDataChannel)
		//
		this.Debugf(configuration.Debug, FormatStruct, configuration)
	}()

	this.Debugf(configuration.Debug,FormatStruct,  configuration)

	for {
		var jobData *JobData
		select {
		case jobData = <-jobDataChannel:
			{
				//check status, if it error then scheduled it again
				if jobData.Error {
					this.Debugf(configuration.Debug,FormatStruct,jobData)
					//increment errors counter
					taskData.Errors++
					//record last dropped
					taskData.LastDroppedJob=*jobData
					if jobData.Attempts == configuration.MaxAttempts {
						//record how many has been dropped- it will be only one
						taskData.MaxAttemptJobsDropped++
						//
						taskData.LastDroppedJob.LastErrorMsg="max attempts reached, reason:"+jobData.LastErrorMsg
						//
						taskData.Failed=true
						//we quit on first dropped job
						return taskData
					}
					//reset flag and try again
					jobData.Error = false

				} else {
					taskData.Success++
					//record the fact that that maany rows were affected by this job
					//configuration.RowsAffected(jobData.RowsAffected)
					taskData.RowsAffected=taskData.RowsAffected+jobData.RowsAffected
				}

				//got report back, decrease the length
				taskData.QueueLength--
			}
		//
		case <-time.After(time.Duration(taskData.Timeout) * time.Millisecond):
		//cancellation
		case <-GDone:
			//return current state of task
			return taskData
		}

		//if queue length reached its limit
		// or exhausted stream
		// or finished
		// then we cannot schedule more jobs
		if (taskData.QueueLength < configuration.Concurrency	&&	((taskData.Success+taskData.MaxAttemptJobsDropped+taskData.QueueLength) < JobsCount)){

			//if current job is empty then we create new job
			if jobData == nil {
				var Query string	=	configuration.Exec
				//if partitioning enabled then format SQL string to provide Min and Max
				if JobsCount>1{
					Query=fmt.Sprintf(configuration.Exec, taskData.DataCursor, taskData.DataCursor + taskData.TaskConfiguration.Step)
				}
				//try to resolve $LastEtl to date time
				Query=strings.Replace(Query,LastEtlVariableName,LastEtl.Format(SessionFileFormat),-1)
				//try to resolve $ToEtl do date time
				Query=strings.Replace(Query,EtlToVariableName,EtlTo.Format(SessionFileFormat),-1)
				//store only data that is requred and specifc for job
				jobData = &JobData{Id: taskData.JobId, Query: Query}
				//move cursor to the next step
				taskData.DataCursor += configuration.Step
				//
				taskData.JobId++
			}

			jobContext:=	JobExecutionContext{
				JobData:jobData,
				Dsn:configuration.Dsn,
				PreSteps:configuration.PreSteps,
				JobDataChannel:jobDataChannel,
				Debug:configuration.Debug,
				TaskName:configuration.Name}
			this.Debugf(configuration.Debug,FormatStruct,jobContext)
			//schedule the job
			//todo: switch case here: Exec, Query, QueryOne
			go this.Exec(jobContext)
			//increase queue length
			taskData.QueueLength++
			//reset timeout
			taskData.Timeout = 0
		} else {
			//if timeout limit is not reached then increase timeout
			if time.Duration(taskData.Timeout) <= time.Second {
				//increase timeout when queue if full or there is no more job to schedule
				taskData.Timeout *= 10
			}
		}
		//have we finished yet?
		if (taskData.Success+taskData.MaxAttemptJobsDropped+taskData.QueueLength) >= JobsCount {
			taskData.Status = "finishing"
			if taskData.QueueLength == 0 {
				taskData.Status = "finished"
				return taskData
			}
		}
		//if last serialisation happend more than...
		if time.Now().Sub(taskData.Serialised) > (1 * time.Second) {
			SerialiseStruct(this.Configuration)
			taskData.Serialised = time.Now()
		}
	}
}

///Atomic SQL exec.
//tested on MySQL update
//todo:handle errors-
//todo:implement query
//todo: implement queryrow - that returns at most one row
//Exec executes a query without returning any rows
func (this *SessionController) Exec(jobContext JobExecutionContext) {
	//how many rows were affected by Exec() (if supported by SQL driver)
	//var totalRowsAffected uint64

	//log some debuf info if in debug mode
	this.Debugf(jobContext.Debug,FormatStruct,jobContext)
	this.Debugf(jobContext.Debug,FormatStruct,jobContext.JobData)

	//store start time- for reporting purposes
	jobContext.JobData.StartTime = time.Now()

	//all deffered functions
	defer func() {
		if err := recover(); err != nil {
			this.Debugf(jobContext.Debug,"panic %s\n",err)
			jobContext.JobData.Error 		= 	true
			jobContext.JobData.LastErrorMsg	=	fmt.Sprintf("%s",err)
		}
		//increase number of attmpts
		jobContext.JobData.Attempts++
		//record data
		jobContext.JobData.StopTime = time.Now()
		//notify producer that another job has finished
		jobContext.JobDataChannel <- jobContext.JobData
		//
		this.Debugf(jobContext.Debug,FormatStruct,jobContext)
		this.Debugf(jobContext.Debug,FormatStruct,jobContext.JobData)
	}()

	//how to use connection pool?
	db, err := sql.Open("mysql", jobContext.Dsn)
	if err != nil {
		this.Debugf(jobContext.Debug,FormatStruct,err)
		panic(err)
	}
	//close connection hence no side effects
	defer db.Close()
	//log.Print("connection open ", Dsn)

	//iterate all 'set'
	for _,stmt:=range jobContext.PreSteps {
		_, err := db.Exec(stmt)
		if err != nil {
			this.Debugf(jobContext.Debug,FormatStruct,err)
			panic(err)
		}
		this.Debugf(jobContext.Debug,"pre-exec: %s",stmt)
	}

	var result sql.Result
	//all data source details should be well encapsulated
	result, err = db.Exec(jobContext.JobData.Query)

	if err != nil {
		this.Debugf(jobContext.Debug,FormatStruct,err)
		panic(err)
	}
	//if driver supports rows affected and last inserted id
	if rowsAffected,err:=result.RowsAffected();err==nil{
		jobContext.JobData.RowsAffected=uint64(rowsAffected)
	}

	this.Debugf(jobContext.Debug,"exec query %s \n",jobContext.JobData.Query)
	this.Debugf(jobContext.Debug,"rows affected: %d\n",jobContext.JobData.RowsAffected)
}


func (this *SessionController)Debugf(enabled bool,format string, args... interface{}){
	if enabled {
		GDLogger.Debugf(format,args)
	}
	return
}
