package main


//comfoguration todo:
//	add storage_type={mysql|postgress|mongodb...}
//	add type of operation:{update|insert...}
//todo error handling:
//	MySQL schema issues

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	_"strings"
	_"strconv"
)

type JobConfiguration struct {
}

const (
	InvalidJob            = 0
	TaskMethod_SQLUpdate  = "SQLUpdate"
	TaskMethod_TestUpdate = "TestUpdate"
	DefaultStatusFileName = "tasks.json"
	//must be exact match
	LastEtlVarName			="$LastEtl"
	MainLiteral ="MAIN:"
)
var(
	//will be used to notify users about some important facts
	Notifier *NotificationsModule
	//global module ued by specialised interfaces for recording session,task, job progress
	GLogger Logger			=	&LoggingModule{}
	//global debug logger
	GDLogger DebugLogger	=	&DebugLoggingModule{}
	//
	GCLogger ConsoleLogger	=	&ConsoleLoggerModule{}

	//current session id
	GSessionId string		=	time.Now().Format(SessionFileFormat)
	//
	GDone					=	make(chan struct{})
	//
	GSessionConfiguration	*SessionConfiguration

	//
	GSessionExecutionContext SessionExecutionContext
)

var (
//last sucessfull etl
//todo: change name to something more generic. move it to session context when available
	LastEtl	time.Time
	LastEtlFlag			=	flag.String("last_etl","","last etl date, not mandatory")
//
	EtlTo time.Time
	EtlToFlag		= flag.String("etl_to","","etl data up to date, not mandatory")
//
	ConfigFileNameFlag	= flag.String("config", "", "configuration file name")
//
	TestConfigLoadFlag	= flag.Bool("test_config",false,"test configuration?")
)

func main() {
	var err error

	//parse parameters
	flag.Parse()

	//if config file name is empty
	if *ConfigFileNameFlag == "" {
		usage()
	}

	//if types dont match then we get zero value for that SessionConguration
	if GSessionConfiguration,err=loadSessionConfiguration(*ConfigFileNameFlag);err != nil {
		Printf("configuration load error: %s",err.Error())
		os.Exit(1)
	}

	Printf("Configuration loaded. Found %d tasks", len(GSessionConfiguration.Tasks))

	//set last, etl to global variable
	//this flag is used onyl for testing purposes
	EtlTo,_		=	time.Parse(LastEtlFileFormat,*EtlToFlag)
	//set last etl to global variable
	LastEtl,_	=	time.Parse(LastEtlFileFormat,*LastEtlFlag)
	//if last etl date was not provided then ww try to find it mongst log files
	if LastEtl.IsZero() {
		//find last_elt timestamp
		if lastEtl, err := findLastEtlTime(); err!=nil {
			Printf("error finding $LastEtl:%s\n", err)
			os.Exit(1)
		}else{
			LastEtl=lastEtl
		}
	}

	Printf("$LastEtl = %s",LastEtl.Format(LastEtlFileFormat))
	Printf("$EtlTo = %s",EtlTo.Format(LastEtlFileFormat))

	if !EtlTo.IsZero()&&EtlTo.Sub(LastEtl)<=0{
		Printf("EtlTo must not be smaller than LastEtl")
		os.Exit(1)
	}

	//do we test config only?
	if *TestConfigLoadFlag{
		Printf("testing config only\n")
		os.Exit(0)
	}

	var sessionController 	*SessionController
	//start new session
	if sessionController,err=	makeSessionController(*GSessionConfiguration);err!=nil{
		Printf("error while creating session log file: %s",err.Error())
		os.Exit(1)
	}

	//this is where fun begins
	sessionController.StartTasks()
	//by closing channel we terminate session
	close(GDone)
}

//
func loadSessionConfiguration(configurationFileName string) (configuration*SessionConfiguration,err error) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			switch v := recovered.(type){
				case error:err=v
			}
			err=StringError(fmt.Sprintf("unknown error, details: %+v", err))
		}else{
			err=nil
		}
	}()

	file,err:=ioutil.ReadFile(configurationFileName)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(file,&configuration); err != nil {
		panic(err)
	}

	if err:= configuration.Init();err!=nil{
		panic(err)
	}

	return configuration,nil
}

func handleError(recoveredValue interface{}) ( error){
	var err	LoggingError
	switch v:=recoveredValue.(type){
		case error: err.Msg=v.Error();err.error=err
		case nil: return nil
	}
	Printf("%s: Error: %s",MainLiteral,err.Error())
	return err
}


func SerialiseStruct(v interface{}) {
	defer func() {
		handleError(recover())
	}()
	bytes, _ := json.Marshal(v)
	//owner=read+wqrite, group and others=read
	ioutil.WriteFile(DefaultStatusFileName, bytes, 0644)
}
func Printf(format string,args... interface{}){
	GCLogger.Printf(MainLiteral+format+"\n",args)
}

func usage(){
	Printf("Data integration tool. Usage:\n")
	Printf("pupdate -config=monitoring/configuration_file [-test_config]\n")
	Printf("-config path to configiration file name\n")
	Printf("-test_config whether to test configuration\n")
	os.Exit(1)
}