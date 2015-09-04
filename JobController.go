package main

import(
	"fmt"
	"time"
)

type JobController struct {
	JobDataChannel chan (*JobData)
}

//
func (this *JobController) StartTask(taskid uint64, configuration *TaskConfiguration) TaskData {
	//create new task data structure where we hold op data and configuration
	var taskData 			= TaskData{Id: taskid, DataCursor: configuration.Min, Status: "working", TaskConfiguration: *configuration}
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
