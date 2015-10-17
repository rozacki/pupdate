package main
import "time"

//as job is not interested in how
// -to communicate with controller
// - how to store results
//- how to trigger error etc
//- where the configuration come from
type JobExecutionContextInterface interface{
	//sets error and finishes then jobi
	SetErrorMessage(msg string)
	GetDsn()string
	GetPreSteps()[]string
	GetPostSteps()[]string
	IsDebug()bool
	TaskName()string
	//read-only hence copy returned or readonly interface
	GetTaskConfiguration() TaskConfiguration
	//
	SetValue(value interface{}) error
	//finishes the job by sending to notyfication to the controller
	Finish()
	//
	StartTime()
	StopTime()
	//
	IncreaseAttempts()
	//
	GetQuery()string
	//
	GetJobData()JobData
	//
	SetRowsAffected(affected uint64)
	GetRowsAffected()uint64
}

type JobExecutionContext struct{
	//what has been done, result, error
	jobData *JobData
	//to send to parent task JobData result with the parent task
	jobDataChannel chan *JobData
	//todo: interfaces data storage: can be memory, file system etc
	//now it is only interface
	value interface{}
	//as contoller has to now details of task to handle response properly
	TaskConfiguration TaskConfiguration
}

//sets error and finishes then jobi
func (this* JobExecutionContext) SetErrorMessage(msg string){
	this.jobData.LastErrorMsg=msg
	this.jobData.Error=true
}
func (this* JobExecutionContext) GetDsn()string{
	return this.TaskConfiguration.Dsn
}

func (this* JobExecutionContext) GetPreSteps()[]string{
	return this.TaskConfiguration.PreSteps
}
func (this* JobExecutionContext) GetPostSteps()[]string{
	return this.TaskConfiguration.PreSteps
}
func (this* JobExecutionContext) IsDebug()bool{
	return this.TaskConfiguration.Debug
}
func (this* JobExecutionContext)  TaskName()string{
	return this.jobData.Name
}
//read-only hence copy returned or readonly interface
func (this* JobExecutionContext) GetTaskConfiguration() TaskConfiguration{
	return this.TaskConfiguration
}
//
func (this* JobExecutionContext) SetValue(value interface{}) error{
	this.value	=	value
	return nil
}
//finishes the job by sending to notyfication to the controller
func (this* JobExecutionContext) Finish(){
	this.jobDataChannel<-this.jobData
}
//
func (this* JobExecutionContext) StartTime(){
	this.jobData.StartTime=time.Now()
}
func (this* JobExecutionContext) StopTime(){
	this.jobData.StopTime	=	time.Now()
}
//
func (this* JobExecutionContext) IncreaseAttempts(){
	this.jobData.Attempts++
}
//
func (this* JobExecutionContext) GetQuery()string{
	return this.jobData.Query
}
//
func (this* JobExecutionContext) GetJobData()JobData{
	return *this.jobData
}
//
func (this* JobExecutionContext) SetRowsAffected(affected uint64){
	this.jobData.RowsAffected=affected
}
func (this* JobExecutionContext) GetRowsAffected()uint64{
	return this.jobData.RowsAffected
}


//TaskConfiguration extands existing MonitoringModule
func (this* JobExecutionContext) EventStartJob()(error){
	return GLogger.Tracef(this.TaskConfiguration.Name,StartJob)
}

func (this* JobExecutionContext) EventStopTJob()(error){
	return GLogger.Tracef(this.TaskConfiguration.Name,StopJob)
}

func (this* JobExecutionContext) Event(data interface{})(error){
	return GLogger.Tracef(this.TaskConfiguration.Name,Trace)
}

//for debugging purposes, does not return error
func (this* JobExecutionContext) Debugf(format string,args... interface{}){
	if !this.TaskConfiguration.Debug{
		return
	}
	GDLogger.Debugf(format,args)
}