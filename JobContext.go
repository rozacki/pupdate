package main

type JobContext struct{
	JobData *JobData
	Dsn string
	PreSteps []string
	JobDataChannel chan *JobData
	//inherited from Task
	Debug bool
	//current data
	TaskName string
}

//TaskConfiguration extands existing MonitoringModule
func (this* JobContext) EventStartJob()(error){
	return GMonitoring.Trace(this.TaskName,StartJob)
}

func (this* JobContext) EventStopTJob()(error){
	return GMonitoring.Trace(this.TaskName,StopJob)
}

func (this* JobContext) Event(data interface{})(error){
	return GMonitoring.Trace(this.TaskName,Trace)
}