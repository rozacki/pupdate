package main

import(
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
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
	if err:= GLogger.OpenLog();err!=nil{
		return nil,err
	}

	return &sessionController,nil
}

//Runs job in order
func (this *SessionController) StartSession() (err error){
	defer func(){
		if recVal:=recover();recVal!=nil{
				switch v:=recVal.(type){
					case error: err=v
					case string: err=StringError(v)
					default: err=StringError("unknown error")
				}
				this.Configuration.SessionFailed(err.Error())
				GLogger.CloseLog(false)
		}else{
			this.Configuration.SessionSuccess()
			GLogger.CloseLog(true)
		}
	}()

	for _,task:=range this.Configuration.Tasks{
		//if task is disabled then we skip it
		if task.Disabled{
			this.Configuration.TaskDisabled(task.Name)
			continue
		}
		//
		this.Configuration.TaskCounter++

		//blocking, synchronous task
		taskData:=this.executeTask(task)

		//on first task failed terminates session
		if taskData.Failed{
			panic(StringError(fmt.Sprintf("'%s'. Task details %+v",taskData.LastDroppedJob.LastErrorMsg,taskData)))
		}

		//should I finish
		select{
		//this channel is used as global flag to finish processing. it is used by main thread and allow user to signal termination
		case <-GDone:
			//todo: if there is at least one more task set status as terminated
			panic("session terminated")
		case  <-time.After(1* time.Millisecond):
		}
	}
	return nil
}
//runs single task by spawning multiple paraller jobs if necessary
//todo: move to taskcontroller
func (this *SessionController) executeTask(taskConfiguration TaskConfiguration) TaskData {

	taskConfiguration.EventStartTask()

	//create new task data structure where we hold op data and configuration
	var taskData 			= 	TaskData{
		Id					: this.Configuration.TaskCounter,
		DataCursor			: taskConfiguration.Min,
		Status				: "working",
		TaskConfiguration	: taskConfiguration,
		CreationTime		: time.Now()}

	//before count jobscount I have to ask mysql for max and min values
	if taskConfiguration.Max==0 && taskConfiguration.MaxQuery!=""{

	}

	if taskConfiguration.Min==taskConfiguration.Max && taskConfiguration.MinQuery!=""{

	}

	var JobsCount	uint64
	if taskConfiguration.Step<=0 || taskConfiguration.Max<=taskConfiguration.Min || taskConfiguration.Concurrency<1{
		JobsCount						=1
		taskConfiguration.Step			=0
		taskConfiguration.Max			=0
		taskConfiguration.Min			=0
		taskConfiguration.Concurrency	=1
	}else {
		JobsCount	= ((taskConfiguration.Max - taskConfiguration.Min)/taskConfiguration.Step)
	}
	//when job finished then posts its status here
	var jobDataChannel		=	make(chan *JobData, taskConfiguration.Concurrency)
	//
	defer func(){
		taskConfiguration.RowsAffected(taskData.RowsAffected)
		if taskData.Failed {
			//record that fact
			taskConfiguration.EventFailTask(taskData.LastDroppedJob.LastErrorMsg)
		}else{
			taskConfiguration.EventSuccessTask()
		}
		//last serialisation
		SerialiseStruct(this.Configuration)
		//close task<-job comm channel
		close(jobDataChannel)
		//
		this.Debugf(taskConfiguration.Debug, FormatStruct, taskConfiguration)
	}()

	this.Debugf(taskConfiguration.Debug,FormatStruct, taskConfiguration)

	for {
		var jobData *JobData
		select {
		case jobData = <-jobDataChannel:
			{
				//check status, if it error then scheduled it again
				if jobData.Error {
					this.Debugf(taskConfiguration.Debug,FormatStruct,jobData)
					//increment errors counter
					taskData.Errors++
					//record last dropped
					taskData.LastDroppedJob=*jobData
					if jobData.Attempts == taskConfiguration.MaxAttempts {
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
		if (taskData.QueueLength < taskConfiguration.Concurrency	&&	((taskData.Success+taskData.MaxAttemptJobsDropped+taskData.QueueLength) < JobsCount)){

			//if current job is empty then we create new job
			if jobData == nil {
				var Query string	=	taskConfiguration.Exec
				//if partitioning enabled then format SQL string to provide Min and Max
				if JobsCount>1{
					Query=fmt.Sprintf(taskConfiguration.Exec, taskData.DataCursor, taskData.DataCursor + taskData.TaskConfiguration.Step)
				}
				//try to resolve $LastEtl to date time
				Query=strings.Replace(Query,LastEtlVariableName,LastEtl.Format(SessionFileFormat),-1)
				//try to resolve $ToEtl do date time
				Query=strings.Replace(Query,EtlToVariableName,EtlTo.Format(SessionFileFormat),-1)
				//store only data that is requred and specifc for job
				jobData = &JobData{Id: taskData.JobId, Query: Query}
				//move cursor to the next step
				taskData.DataCursor += taskConfiguration.Step
				//
				taskData.JobId++
			}

			jobContext:=	&JobExecutionContext{
				jobData:jobData,
				jobDataChannel:jobDataChannel}
			this.Debugf(taskConfiguration.Debug,FormatStruct,jobContext)
			//schedule the job
			//todo: switch case here: Exec, Query, QueryOne
			switch strings.ToUpper(taskConfiguration.ExecType){
				case Exec:	go this.Exec(jobContext)
				case Query: go this.QueryRow(jobContext)
				case QueryOne:
			}

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
func (this *SessionController)Debugf(enabled bool,format string, args... interface{}){
	if enabled {
		GDLogger.Debugf(format,args)
	}
	return
}
