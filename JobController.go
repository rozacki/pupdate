package main

import(
	"fmt"
	"time"
)

type SessionController struct {
	SessionID	string
	JobDataChannel chan (*JobData)
	TaskCounter	uint64
	Configuration SessionConfiguration
}

func makeSession(configuration SessionConfiguration)(session* SessionController){
	//is there  any task to schedul?
	if len(configuration.Tasks)==0{
		fmt.Println("no tasks to schedule, stop")
		return nil
	}

	var sessionController=	SessionController{SessionID:time.Now().Format(time.UnixDate),Configuration:configuration}
	sessionController.JobDataChannel=make(chan *JobData, configuration.Tasks[0].Concurrency)
	return &sessionController
}

//Runs job in order
func (this *SessionController) StartTasks(configuration *TaskConfiguration) {
	for task:=range this.Configuration.Tasks{
		taskData:=this.startTask(task)
		if taskData.Success==false{
			fmt.Printf("task failed:%n",)
		}
	}
}
//
func (this *SessionController) startTask(configuration *TaskConfiguration) TaskData {
	//
	this.TaskCounter++
	//create new task data structure where we hold op data and configuration
	var taskData 			= 	TaskData{Id: this.TaskCounter, DataCursor: configuration.Min, Status: "working", TaskConfiguration: *configuration}
	var NoJobs				=	((configuration.Max - configuration.Min)/configuration.Step)

	if configuration.Debug {
		fmt.Printf("%+v\n", configuration)
	}

	for {
		var jobData *JobData
		select {
		case jobData = <-this.JobDataChannel:
			{
				//check status, if it error then scheduled it again
				if jobData.Error {
					if configuration.Debug {
						fmt.Printf("%+v\n", jobData)
					}
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
		}

		//if queue length reached its limit
		// or exhausted stream
		// or finished
		// then we cannot schedule more jobs
		if taskData.QueueLength < configuration.Concurrency&& (taskData.Success+taskData.MaxAttemptJobsDropped+taskData.QueueLength) < NoJobs{

			//if current job is empty then we create new job
			if jobData == nil {
				//store only data that is requred and specifc for job
				jobData = &JobData{Id: taskData.JobId, PartStart: taskData.DataCursor, PartEnd: taskData.DataCursor + taskData.TaskConfiguration.Step}
				//move cursor to the next step
				taskData.DataCursor += configuration.Step
				//
				taskData.JobId++
			}
			//todo: switch case here
			//schedule the job
			go this.SQLUpdate(jobData, configuration.Dsn, configuration.SessionParam, fmt.Sprintf(configuration.Update, jobData.PartStart, jobData.PartEnd))
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
		if (taskData.Success+taskData.MaxAttemptJobsDropped+taskData.QueueLength) >= NoJobs {
			taskData.Status = "finishing"
			if taskData.QueueLength == 0 {
				taskData.Status = "finished"
				SerialiseStruct(taskData)
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


//tested on MySQL updated
func (this *SessionController) SQLUpdate(jobData *JobData, dsn string, sessionParams string, query string) {
	//store start time
	jobData.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			//log.Println(err)
			jobData.Error = true
			//jobData.ErrorMsg=err.Error()
		}
		//increase number of attmpts
		jobData.Attempts++
		//record data
		jobData.StopTime = time.Now()
		//notify producer that another job has finished
		this.JobDataChannel <- jobData
	}()

	//fmt.Printf("start: job_id=%d, start_id=%d stop_id=%d\n",jobId, startid, limit)

	//how to use connection pool?
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//log.Print("connection open ", Dsn)

	if len(sessionParams) > 0 {
		_, err = db.Exec(sessionParams)
	}

	if err != nil {
		panic(err)
	}

	//all data source details should be well encapsulated
	_, err = db.Exec(query)

	//log.Print("query finished:", query)

	if err != nil {
		panic(err)
	}

	//fmt.Printf("end: job_id=%d,start_id=%d stop_id=%d\n",jobId, startid, limit)
}

