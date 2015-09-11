package main

import (
	"fmt"
	"os"
	"path/filepath"
	_"time"
	_"path"
	_"sort"
)

const(
	StartSession="StartSession"
	StartTask="StartTask"
	StopSession="StopSession"
	StopTask="StopTask"
	Trace="Trace"
	StartJob="StartJob"
	StopJob="StopJob"

	SessionLogFileExt 		=	".log"
	DatSessionLogFileExt =	".dat"
	SessionLogFolder		=	"monitoring"

	//this is how it is stored as session file name, up to seconds
	SessionFileFormat	  = "2006-01-02 15:04:05"
	//this is how it is presented to SQL, daily granularity
	LastEtlFileFormat	  =	"2006-01-02"
	//last etl variable name used in sql staements
	LastEtlVariableName		="$last_etl"
)

//Monitoring module
//todo:
type MonitoringModule struct{
	Configuration MonitoringConfiguration
	File *os.File
	FileName string
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

func (this*MonitoringModule) Trace(sid string,taskName string,tid uint64,jid uint64,event string,data interface{})(*MonitoringError){
	return this.TraceOK(sid,taskName ,tid ,jid ,event ,data,false)
}

func (this*MonitoringModule) TraceOK(sid string,taskName string,tid uint64,jid uint64,event string,data interface{}, ok bool)(*MonitoringError){
	if len(sid)==0{
		this.log("missing session id")
		return nil
	}
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
	var err error
	//is log file specific for this sesion is open now?
	//I don't check if event type is "NewSession"
	if this.File==nil{
		//if file sid file is not open and does not exist then open it
		//if file is not open and does exist then fmt.Println() and return
		this.FileName=filepath.Join(SessionLogFolder,sid+SessionLogFileExt)
		if this.File,err=os.OpenFile(this.FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY| os.O_EXCL, 0600);err!=nil{
			this.log("error '%s' during creating session file '%s' ",err.Error(), this.FileName)
			return nil
		}
		this.log("session file '%s' created",this.FileName)
	}

	s:= fmt.Sprintf("\n %+v \n",Event)
	if _, err = this.File.WriteString(s); err != nil {
		return nil
	}

	//if event=="StopSession" then write event, close the file and change the name
	if event==StopSession{
		this.File.Close()
		this.log("session log closed")
		if ok{
			os.Rename(this.FileName,filepath.Join(SessionLogFolder,sid+DatSessionLogFileExt))
			this.log("session log renamed")
		}
	}

	return nil
}
func (this*MonitoringModule)log(format string,args... string){
	fmt.Printf("monitoring:"+format+"\n",args)
}
//Depeneding on configuration we can support db monitoring.
// Currently we suport only log-based monitoring
type MonitoringConfiguration struct{
	Dsn string
}

func makeMonitoring(configuration MonitoringConfiguration)(Monitor){
	if len(configuration.Dsn)==0{
		fmt.Printf("missing dsn for monitoring module")
		return nil;
	}
	//we support only one monitoring type
	return &MonitoringModule{Configuration:configuration}
}

//specific interface for monitoring tasks
type Monitor interface{
	//generic method
	Trace(sid string,taskName string,tid uint64,jid uint64,event string,data interface{})(*MonitoringError)
	//generic method thta accespts success ot false
	TraceOK(sid string,taskName string,tid uint64,jid uint64,event string,data interface{}, ok bool)(*MonitoringError)
}
/*
func findLastEtlTime() time.Time{
	if matches, err:= filepath.Glob(path.Join(SessionLogFolder,"*"+DatSessionLogFileExt));err!=nil{
		return nil
	}

	//sort
	sortedFileNames:=sort.Strings(ByLength(matches))

	time.Parse(SessionFileFormat)
}
type ByLength []string
func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}
*/