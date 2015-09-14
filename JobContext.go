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
}

//TaskConfiguration extands existing MonitoringModule
func (this* JobContext) EventStartJob()(*MonitoringError){
	return Monitoring.Trace("","",0,this.JobData.Id,StartJob,this.JobData)
}

func (this* JobContext) EventStopTJob()(*MonitoringError){
	return Monitoring.Trace("","",0,this.JobData.Id,StopJob,this.JobData)
}

func (this* JobContext) Event(data interface{})(*MonitoringError){
	return Monitoring.Trace("","",0,this.JobData.Id,Trace,data)
}