package main

type JobContext struct{
	JobData *JobData
	Dsn string
	PreSteps []string
	JobDataChannel chan *JobData
	//inherited from Task
	Debug bool
	//how many rows have been affected based on what driver returns
	RowsAffected uint64
	//current data
	TaskName string
}

//TaskConfiguration extands existing MonitoringModule
func (this* JobContext) EventStartJob()(*MonitoringError){
	return Monitoring.Trace(this.TaskName,StartJob)
}

func (this* JobContext) EventStopTJob()(*MonitoringError){
	return Monitoring.Trace(this.TaskName,StopJob)
}

func (this* JobContext) Event(data interface{})(*MonitoringError){
	return Monitoring.Trace(this.TaskName,Trace)
}