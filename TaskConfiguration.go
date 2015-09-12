package main

import(

	"strings"
	_"fmt"
)

type TaskConfiguration struct {

	Name   string
	Dsn    string
	Exec   string
	ExecTab[] string
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
	TaskId uint64
}
//TaskConfiguration extands existing MonitoringModule
func (this* TaskConfiguration) EventStartTask()(*MonitoringError){
	return Monitoring.Trace(this.SessionId,this.Name,this.TaskId,0,StartTask,this)
}

func (this* TaskConfiguration) EventStopTask()(*MonitoringError){
	return Monitoring.Trace(this.SessionId,this.Name,this.TaskId,0,StopTask,this)
}

func (this* TaskConfiguration) Trace(data interface{})(*MonitoringError){
	return Monitoring.Trace(this.SessionId,this.Name,this.TaskId,0, Trace,data)
}
//does some housekeeping int he task configuration, shoudl be called after configuration is loaded
func (this* TaskConfiguration) Init() error{
	//join ExecTab into Exec if Exec is emopty and ExecTab is not
	this.Exec=strings.Join(this.ExecTab,"")
	//fmt.Println(this.Exec)
	return nil
}