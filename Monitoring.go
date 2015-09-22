package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"path"
	_"sort"
	"strings"
)

const(
	StartSession="StartSession"
	StartTask="StartTask"
	StopSession="StopSession"
	StopTask="StopTask"
	Trace="Trace"
	StartJob="StartJob"
	StopJob="StopJob"
	SessionSuccess	=	"SessionSuccess"
	SessionFail		=	"SessionFail"

	SessionLogFileExt 		=	".log"
	DatSessionLogFileExt =	".dat"
	SessionLogFolder		=	"monitoring"

	//this is how it is stored as session file name, up to seconds
	SessionFileFormat	  = "2006-01-02 15:04:05"
	//this is how it is presented to SQL, daily granularity
	LastEtlFileFormat	  =	"2006-01-02"
	//last etl variable name used in sql staements
	//todo: change name into last_success_session
	LastEtlVariableName		="$last_etl"
	EtlToVariableName		= "$etl_to"
)

//Monitoring module
//todo:use standard logger. and its API. Std logger gives you a way to controll file handler. Srd logger is thread safe (thanks to io.Writer?)
//automatically call notifier for critical issues or starting and stopping session
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

func (this*MonitoringModule) Trace(taskName string,event string)(*MonitoringError){
	return this.TraceOK(taskName,event,false)
}

func (this*MonitoringModule) TraceOK(taskName string,event string, ok bool)(*MonitoringError){

	Ev:= fmt.Sprintf("event:%s, task name:%s",event,taskName)
	var err error
	//is log file specific for this sesion is open now?
	//I don't check if event type is "NewSession"
	if this.File==nil{
		//if file sid file is not open and does not exist then open it
		//if file is not open and does exist then fmt.Println() and return
		this.FileName=filepath.Join(SessionLogFolder,GSessionId+SessionLogFileExt)
		if this.File,err=os.OpenFile(this.FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY| os.O_EXCL, 0600);err!=nil{
			this.Printf("error '%s' during creating session file '%s' ",err.Error(), this.FileName)
			return nil
		}
		this.Printf("session file '%s' created",this.FileName)
	}

	s:= fmt.Sprintf("%s \n", Ev)
	if _, err = this.File.WriteString(s); err != nil {
		return nil
	}

	//if event=="StopSession" then write event, close the file and change the name
	if event==StopSession || event==SessionSuccess{
		this.File.Close()
		this.Printf("session log closed")
		if ok{
			os.Rename(this.FileName,filepath.Join(SessionLogFolder,GSessionId+DatSessionLogFileExt))
			this.Printf("session log renamed")
		}
	}

	return nil
}
func (this*MonitoringModule)Printf(format string,args... string){
	if len(args)==0{
		fmt.Print(format)
	}else {
		fmt.Printf("monitoring:"+format+"\n",args)
	}
}
// Depending on configuration we can support db monitoring.
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

//specific interface for monitoring sessions, tasks, jobs. It wraps logger api into a number of domain specific methods.
type Monitor interface{
	//generic method
	Trace(taskName string,event string,)(*MonitoringError)
	//generic method thta accespts success ot false
	TraceOK(taskName string,event string, ok bool)(*MonitoringError)
}

func findLastEtlTime() (time.Time,error) {
	var currTime time.Time
	var matches []string
	var err error
	if matches, err= filepath.Glob(path.Join(SessionLogFolder, "*"+DatSessionLogFileExt)); err!=nil {
		return currTime, err
	}

	for _,filePath := range matches {
		_,fileName:=filepath.Split(filePath)
		//strip extension
		fileName=strings.TrimSuffix(fileName,filepath.Ext(fileName))
		//parse the name of the file according to SessionFileFormat
		time, err := time.Parse(SessionFileFormat, fileName)
		if err!=nil {
			continue
		}
		if currTime.Sub(time)<0{
			currTime=time
		}
	}
	return currTime,nil
}