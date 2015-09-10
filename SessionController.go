package main

import(
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type SessionController struct {
	Configuration SessionConfiguration
}

func makeSession(configuration SessionConfiguration)(session* SessionController){
	//is there  any task to schedule?
	if len(configuration.Tasks)==0{
		fmt.Println("no tasks to schedule, stop")
		return nil
	}

	//generate session name
	sessionController:=SessionController{Configuration:configuration}
	//try to record progress, if fails it means that monitoring module is not ok
	if err:=configuration.StartSession();err!=nil{
		return nil
	}
	return &sessionController
}

//Runs job in order
func (this *SessionController) StartTasks() {
	for _,task:=range this.Configuration.Tasks{
		//
		this.Configuration.TaskCounter++
		//
		task.SessionId	=	this.Configuration.SessionID
		//synchronous task
		this.startTask(task,this.Configuration.Done)
		//
		//Monitoring
		//todo:check success
		//if sucess then store task.CreationTime and flag success
		/*
		if taskData.Success==false{
			fmt.Printf("task failed:%n",)
		}
		*/
		//should I finish
		select{
		//finished when task reacts on channel clousure
		case <-this.Configuration.Done:
			fmt.Printf("task exiting\n")
			return
		case  <-time.After(1* time.Millisecond):
		}
	}
	this.Configuration.StopSession()
}
//runs single task by spawning jobs
//todo: cancellation via done channel
//todo: move to taskcontroller
func (this *SessionController) startTask(configuration TaskConfiguration, done <-chan struct{}) TaskData {
	configuration.EventStartTask()
	defer func(){
		configuration.EventStopTask()
	}()
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
	var jobDataChannel		=	make(chan *JobData, configuration.Concurrency)

	ShowDebuginfo(configuration.Debug, configuration)

	for {
		var jobData *JobData
		select {
		case jobData = <-jobDataChannel:
			{
				//check status, if it error then scheduled it again
				if jobData.Error {
					ShowDebuginfo(configuration.Debug,jobData)
					//reset flag and try again
					jobData.Error = false
					//increment errors counter
					taskData.Errors++
					if jobData.Attempts == configuration.MaxAttempts {
						taskData.MaxAttemptJobsDropped++
						jobData = nil
					}
				} else {
					taskData.Success++
				}

				//got report back, decrease the length
				taskData.QueueLength--
			}
		//
		case <-time.After(time.Duration(taskData.Timeout) * time.Millisecond):
		//cancellation
		case <-done:
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
				//if partitioning enabled
				if JobsCount>1{
					Query=fmt.Sprintf(configuration.Exec, taskData.DataCursor, taskData.DataCursor + taskData.TaskConfiguration.Step)
				}
				//store only data that is requred and specifc for job
				jobData = &JobData{Id: taskData.JobId, Query: Query}
				//move cursor to the next step
				taskData.DataCursor += configuration.Step
				//
				taskData.JobId++
			}

			configuration.Event("I'm here")

			jobContext:=	JobContext{JobData:jobData,
				Dsn:configuration.Dsn,
				SessionParams:configuration.SessionParam,
				JobDataChannel:jobDataChannel,
				Debug:configuration.Debug}
			ShowDebuginfo(configuration.Debug,jobContext)
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
				SerialiseStruct(taskData)
				close(jobDataChannel)
				return taskData
			}
		}
		//if last serialisation happend more than...
		if time.Now().Sub(taskData.Serialised) > (1 * time.Second) {
			SerialiseStruct(taskData)
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
func (this *SessionController) Exec(jobContext JobContext) {
	ShowDebuginfo(jobContext.Debug,jobContext)
	ShowDebuginfo(jobContext.Debug,jobContext.JobData)
	//store start time
	jobContext.JobData.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			//log.Println(err)
			jobContext.JobData.Error = true
			//jobData.ErrorMsg=err.Error()
		}
		//increase number of attmpts
		jobContext.JobData.Attempts++
		//record data
		jobContext.JobData.StopTime = time.Now()
		//notify producer that another job has finished
		jobContext.JobDataChannel <- jobContext.JobData
		//
		ShowDebuginfo(jobContext.Debug,jobContext)
		ShowDebuginfo(jobContext.Debug,jobContext.JobData)
	}()

	//how to use connection pool?
	db, err := sql.Open("mysql", jobContext.Dsn)
	if err != nil {
		ShowDebuginfo(jobContext.Debug,err)
		panic(err)
	}
	defer db.Close()
	//log.Print("connection open ", Dsn)

	if len(jobContext.SessionParams) > 0 {
		_, err = db.Exec(jobContext.SessionParams)
	}

	if err != nil {
		ShowDebuginfo(jobContext.Debug,err)
		panic(err)
	}

	fmt.Println(jobContext.JobData.Query)
	//all data source details should be well encapsulated
	_, err = db.Exec(jobContext.JobData.Query)

	//log.Print("query finished:", query)

	if err != nil {
		ShowDebuginfo(jobContext.Debug,err)
		panic(err)
	}

	//fmt.Printf("end: job_id=%d,start_id=%d stop_id=%d\n",jobId, startid, limit)
}

func ShowDebuginfo(b bool, v interface{}){
	defer func(){
		recover()
	}()
	if b{
		fmt.Printf("%#v\n",v)
	}
}