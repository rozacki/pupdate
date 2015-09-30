package main

type JobExecutionContext struct{
	//what has been done, result, error
	JobData *JobData
	//connection string
	Dsn string
	//what steps to take before executing actua job
	PreSteps []string
	//to send to parent task JobData result with the parent task
	JobDataChannel chan *JobData
	//inherited from Task, are we in debug mode
	Debug bool
	//parent task name
	TaskName string
}

//TaskConfiguration extands existing MonitoringModule
func (this* JobExecutionContext) EventStartJob()(error){
	return GLogging.Tracef(this.TaskName,StartJob)
}

func (this* JobExecutionContext) EventStopTJob()(error){
	return GLogging.Tracef(this.TaskName,StopJob)
}

func (this* JobExecutionContext) Event(data interface{})(error){
	return GLogging.Tracef(this.TaskName,Trace)
}

//for debugging purposes, does not return error
func (this* JobExecutionContext) Debugf(format string,args... interface{}){
	if !this.Debug{
		return
	}
	GDLogger.Debugf(format,args)
}