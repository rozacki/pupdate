package main

import (
	"fmt"
)

const(
	StartSession ="session_start"
	StartTask ="new_task"
	StopSession="stop_session"
	StopTask="stop_task"
	Event="event"
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

func (this*MonitoringModule) Event(sid string,taskName string,tid uint64,jid uint64,event string,data interface{})(*MonitoringError){
	Event:=struct{
		Ev string
		SID string
		TaskName string
		TID uint64
		JID uint64
		Data interface{}
	}{
		event,sid ,taskName ,tid ,jid,data,
	}
	//fmt.Printf("Event: %s %s %d %d %s %#v\n",sid,taskName,tid,jid,event,data)
	fmt.Printf("%+v\n",Event)
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
	Event(sid string,taskName string,tid uint64,jid uint64,event string,data interface{})(*MonitoringError)
}