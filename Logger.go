//Logger logs events to session log file. If the session is successful it changes the extension.
package main

import(
	"fmt"
	"path/filepath"
	"os"
	"time"
	"path"
	"strings"
)

const(
	StartSession		=	"StartSession"
	StartTask			=	"StartTask"
	StopSession			=	"StopSession"
	StopTask			=	"StopTask"
	Trace				=	"Trace"
	StartJob			=	"StartJob"
	StopJob				=	"StopJob"
	SessionSuccess		=	"SessionSuccess"
	SessionFailed 		=	"SessionFail"
	RowsAffected		=	"RowsAffected"
	TotalRowsAffected	=	"TotalRowsAffected"
	TaskFailed			=	"TaskFailed"
	TaskSuccess			=	"TaskSuccess"
	TaskDisabled		=	"TaskDisabled"

	SessionLogFileExt 		=	".log"
	DatSessionLogFileExt 	=	".dat"
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

//Generic interface for monitoring sessions, tasks, jobs.
type Logger interface{
	//
	OpenLog()(err error)
	//
	CloseLog(bool)(err error)

	//Methods that output to the log file
	Tracef(format string, args... interface{})(error)
}

//Monitoring module
//todo:use standard logger. and its API. Std logger gives you a way to controll file handler. Std logger is thread safe (thanks to io.Writer?)
type LoggingModule struct{
	File *os.File
	FileName string
}

//handles all write log errors
func (this*LoggingModule)handleLogError(recoveredValue interface{})(error){
	const (LogginError ="Logging error")
	var err LoggingError
	switch v:=recoveredValue.(type){
		case error:
		err.Msg=fmt.Sprintf("%s: %s",LogginError, v.Error())
		err.error=v
		case string:
		err.Msg=fmt.Sprintf("%s: %s",LogginError, v)
		default:
		err.Msg=fmt.Sprintf("%s: %+v",LogginError, recoveredValue)
		case nil:return nil
	}
	this.Printf(err.Msg)
	return err
}

//outputs to the session log file
func (this*LoggingModule) Tracef(format string, args... interface{})(err error){
	defer func(){
		err=this.handleLogError(recover())
	}()
	//output whatever it is
	s:=format
	//if there are some arguments then format
	if len(args)!=0{
		args=PrependArrayOfInterfaces(args,[]interface{}{time.Now().Format(time.StampMilli)})
		s=fmt.Sprintf("%s "+format+"\n",args...)
	}

	if _, err := this.File.WriteString(s); err != nil {
		panic(err)
	}

	return nil
}

func (this* LoggingModule)OpenLog()(err error){
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
func (this*LoggingModule) CloseLog(success bool)(err error){
	defer func(){
		err=this.handleLogError(recover())
	}()
	this.File.Close()
	this.Debugf("session log closed")
	if success{
		os.Rename(this.FileName,filepath.Join(SessionLogFolder,GSessionId+DatSessionLogFileExt))
		this.Debugf("session log renamed")
	}
	return nil;
}

//As logger uses filesystem we need to support additional loggin in case of fs issues. Outputs additional details to the console if Debug=true
func (this*LoggingModule)Debugf(format string,args... interface{}){
	if GSessionConfiguration.Logging.Debug {
		GDLogger.Debugf(format, args)
	}
}
////As logger uses filesystem we need to support additional loggin in case of fs issues. Outputs to console regardless Debug flag
func (this*LoggingModule)Printf(format string,args... interface{}){
	const (Logging	="LOGGING:")
	GCLogger.Printf(Logging+format+"\n", args)
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