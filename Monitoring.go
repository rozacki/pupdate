package main

import (
	"fmt"
	"encoding/json"
)

const(
	StartSession ="session_start"
	StartTask ="new_task"
	StopSession="stop_session"
	StopTask="stop_task"
)

//Monitoring module
//todo:
type MonitoringModule struct{
	Configuration MonitoringConfiguration
}

//Error information sepcific to Monitoring module
type MonitoringError struct{
	Msg string
	Code string
}

//The interface- how not be so verbose? https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
func(this*MonitoringModule) StartSession(sessionid string)(err error){
	return	this.record(fmt.Sprintf("session:%s, event: %s\n",sessionid, StartSession))
}
func(this*MonitoringModule) StopSession(sessionid string)(err error){
	return	this.record(fmt.Sprintf("session:%s, event: %s\n",sessionid,StopSession))
}
func(this*MonitoringModule) StartTask(sessionid string, task string )(err error){
	return	this.record(fmt.Sprintf("session:%s, event: %s\n",sessionid, StartTask))
}
func(this*MonitoringModule) StopTask(sessionid string, task string)(err error){
	return	this.record(fmt.Sprintf("session:%s, event: %s\n",sessionid,StopTask))
}
func(this*MonitoringModule) StartJob(sessionid string, task string, job string)(err error){
	return	this.record(fmt.Sprintf("session:%s, event: %s\n",sessionid, StartTask))
}
func(this*MonitoringModule) StopJob(sessionid string, task string, job string)(err error){
	return	this.record(fmt.Sprintf("session:%s, event: %s\n",sessionid,StopTask))
}

//abstract method that hides storage details
func (this*MonitoringModule) record(msg string)(err error){
	fmt.Println(msg)
	return nil
}

func (this*MonitoringModule) Event(sid string,tid string,jid string,event string,data interface{})(*MonitoringError){
	b,_:=json.Marshal(data)
	fmt.Printf("***%s %s %s %s %s\n",sid,tid,jid,event,string(b))
	return nil
}

type MonitoringConfiguration struct{
	Dsn string
}

func makeMonitoring(configuration MonitoringConfiguration)(Monitor){
	if len(configuration.Dsn)==0{
		fmt.Printf("missing dsn for monitoring module")
		return nil;
	}
	return &MonitoringModule{Configuration:configuration}
}

//specific interface for monitoring tasks
type Monitor interface{
	//generic method
	Event(sid string,tid string,jid string,event string,data interface{})(*MonitoringError)
}