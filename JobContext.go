package main

type JobContext struct{
	JobData *JobData
	Dsn string
	SessionParams string
	JobDataChannel chan *JobData
	Debug bool
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