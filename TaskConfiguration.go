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
	PreSteps []string
	//one of many possible parameters of session
	SessionParamTab []string
	//name of method to run. currenlty supported: SQLUpdate,TestSQL. Default is SQLUpdate
	Method string
	//How many attempts before job is dropped
	MaxAttempts uint64
	//
	Disabled bool
	//Allows to debug individual tasks. This flag is inherited by JonContext.
	Debug bool
//**************** dynamic parameters
	//current session id
	SessionId string
	//current task Id
	TaskId uint64
}
//TaskConfiguration extands existing MonitoringModule
func (this* TaskConfiguration) EventStartTask()(*MonitoringError){
	return Monitoring.Trace(this.Name,StartTask)
}

func (this* TaskConfiguration) EventStopTask()(*MonitoringError){
	return Monitoring.Trace(this.Name,StopTask)
}

func (this* TaskConfiguration) Trace(data interface{})(*MonitoringError){
	return Monitoring.Trace(this.Name,Trace)
}
//does some housekeeping int he task configuration, shoudl be called after configuration is loaded
func (this* TaskConfiguration) Init() error{
	if len(this.ExecTab)>0{
		//join ExecTab into Exec if Exec is emopty and ExecTab is not
		this.Exec=strings.Join(this.ExecTab,"")
		//fmt.Println(this.Exec)
	}
	return nil
}