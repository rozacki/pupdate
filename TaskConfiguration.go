package main

import(

)

type TaskConfiguration struct {

	Name   string
	Dsn    string
	Update string
	//it will be first parameter in the query
	//todo: add special value (string) to indicate that all rows. staring from MIN should be transferred
	Max uint64
	//it will be the second parameter in the query
	Min         uint64
	Step        uint64
	Concurrency uint64
	//one of many possible parameters of session
	SessionParam string
	//name of method to run. currenlty supported: SQLUpdate,TestSQL. Default is SQLUpdate
	Method string
	//How many attempts before job is dropped
	MaxAttempts uint64
	//
	Disabled bool
	//
	Debug bool
//**************** dynamic parameters
	//current session id
	SessionId string
	//current task Id
	TaskId string
}
//TaskConfiguration extands existing MonitoringModule
func (this* TaskConfiguration) EventStartTask()(*MonitoringError){
	return Monitoring.Event(this.SessionId,this.TaskId,"",StartSession,this)
}

func (this* TaskConfiguration) EventStopTask()(*MonitoringError){
	return Monitoring.Event(this.SessionId,this.TaskId,"",StopTask,this)
}