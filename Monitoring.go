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
	SessionFailed =	"SessionFail"
	RowsAffected	=	"RowsAffected"
	TotalRowsAffected	=	"TotalRowsAffected"
	TaskFailed		=	"TaskFailed"
	TaskSuccess		=	"TaskSuccess"
	TaskDisabled	=	"TaskDisabled"

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

	LoggingError		="logging error"
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
type MonitoringError struct {
	//original error
	error
	Msg      string
}

func (this MonitoringError) Error() string{
	return this.Msg
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

func (this*MonitoringModule) Trace(taskName string,event string)(error){
	return this.TraceOK(taskName,event,false)
}

func (this*MonitoringModule) TraceOK(taskName string,event string, ok bool)(err error){
	defer func(){
		err=this.handleLogError(recover())
	}()

	Ev:= fmt.Sprintf("%s, task name:%s",taskName, event)
	s:= fmt.Sprintf("%s %s \n",time.Now().Format(time.StampMilli), Ev)
	if _, err := this.File.WriteString(s); err != nil {
		panic(err)
	}

	return nil
}

func (this*MonitoringModule)handleLogError(recoveredValue interface{})(err error){
	const (LogginError ="Logging error")
	if recoveredValue!=nil {
		var errorMsg string
		switch v:=recoveredValue.(type){
			case error:
			errorMsg=fmt.Sprintf("%s: %s",LogginError, v.Error())
			case string:
			errorMsg=fmt.Sprintf("%s: %s",LogginError, v)
			default:
			errorMsg=fmt.Sprintf("%s: %+v",LogginError, recoveredValue)
		}
		this.Printf(errorMsg)
		return &MonitoringError{Msg:errorMsg}
	}else{
		//will calling function return nil
		return nil
	}
}

func (this* MonitoringModule)OpenLog()(err error){
	defer func(){
		err=this.handleLogError(recover())
	}()

	//is log file specific for this sesion is open now?
	//I don't check if event type is "NewSession"
	if this.File==nil{
		//if file sid file is not open and does not exist then open it
		//if file is not open and does exist then fmt.Println() and return
		this.FileName=filepath.Join(SessionLogFolder,GSessionId+SessionLogFileExt)
		var openError error
		if this.File, openError=os.OpenFile(this.FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY| os.O_EXCL, 0600); openError!=nil{
			this.Debugf("error '%s' during creating session file '%s' ", openError.Error(), this.FileName)
			panic(openError)
		}
		//if we output session name then user will know which file to look into
		this.Printf("session file created '%s' ",this.FileName)
	}
	return nil
}

//Explicit closing log file. If success=true then file's extension is changed.
//If success then file's extension stays intact
func (this*MonitoringModule) CloseLog(success bool)(err error){
	defer func(){
		if err:=recover();err!=nil {
			Printf(LoggingError)
		}
	}()
	this.File.Close()
	this.Debugf("session log closed")
	if success{
		os.Rename(this.FileName,filepath.Join(SessionLogFolder,GSessionId+DatSessionLogFileExt))
		this.Debugf("session log renamed")
	}
	return nil;
}
//outputs additional details to the console if Debug=true
func (this*MonitoringModule)Debugf(format string,args... string){
	if this.Configuration.Debug {
		Printf(format,args)
	}
}
//outputs to console regardless Debug flag
func (this*MonitoringModule)Printf(format string,args... string){
	const (Monitoring	="MONITORING:")
	if len(args)==0 {
		fmt.Println(Monitoring+format)
	}else {
		fmt.Printf(Monitoring+format+"\n", args)
	}
}

//outputs to the session log file
func (this*MonitoringModule)Tracef(taskName string,format string, args... string)(error){
	return this.Trace(taskName,fmt.Sprintf(format,args))
}

// Depending on configuration we can support db monitoring.
// Currently we suport only log-based monitoring
type MonitoringConfiguration struct{
	Dsn string
	Debug bool
}

func makeMonitoring(configuration MonitoringConfiguration)(Monitor,error){
	if len(configuration.Dsn)==0{
		return nil,StringError{"missing dsn for monitoring module"}
	}
	//we support only one monitoring type
	return &MonitoringModule{Configuration:configuration},nil
}

//specific interface for monitoring sessions, tasks, jobs. It wraps logger api into a number of domain specific methods.
type Monitor interface{
	//generic method
	Trace(taskName string,event string)(error)
	//generic method thta accespts success ot false
	TraceOK(taskName string,event string, ok bool)(error)
	//
	Tracef(taskName string,format string, args... string)(error)
	//
	OpenLog()(err error)
	//
	CloseLog(bool)(err error)
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