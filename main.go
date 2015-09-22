package main


//comfoguration todo:
// 	multiple tasks from one configuration
//	add storage_type={mysql|postgress|mongodb...}
//	add type of operation:{update|insert...}
//todo error handling:
//	MySQL schema issues

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	_"strings"
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
)
var(
	//will be used to notify users about some important facts
	Notifier *NotificationsModule
	//global module for recording session,task, job progress
	Monitoring Monitor
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
	//var TestConfigLoadFlag bool
	//current session id
	GSessionId string
)

func main() {
	//parse parameters
	flag.Parse()

	if *ConfigFileNameFlag == "" {
		usage()
	}

	//if file is corrupted then we get nil
	//if types dont match then we get zero value for that SessionConguration
	configuration := loadSessionConfiguration(*ConfigFileNameFlag)
	if configuration == nil {
		fmt.Println("configuration load error")
		os.Exit(1)
	}
	Printf("Configuration loaded. Found %d tasks\n", len(configuration.Tasks))

	if Monitoring=makeMonitoring(configuration.Monitoring);Monitoring==nil{
		os.Exit(1)
	}

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

	Printf("$LastEtl = %s\n",LastEtl.Format(LastEtlFileFormat))
	Printf("$EtlTo = %s\n",EtlTo.Format(LastEtlFileFormat))

	if !EtlTo.IsZero()&&EtlTo.Sub(LastEtl)<=0{
		Printf("EtlTo must not be smaller than LastEtl")
		os.Exit(1)
	}

	GSessionId				=	time.Now().Format(SessionFileFormat)
	configuration.SessionID	=	GSessionId
	configuration.Done		=	make(chan struct{})

	//do we test config only?
	if *TestConfigLoadFlag{
		fmt.Printf("testing config only\n")
		os.Exit(0)
	}

	var sessionController 	*SessionController
	//start new session
	if sessionController=	makeSessionController(*configuration);sessionController==nil{
		os.Exit(1)
	}

	sessionController.StartTasks()
	//it will terminate session
	close(configuration.Done)

	/*
	var key string
	fmt.Println("Press key if you want to finish")
	fmt.Scanf("%s",&key)
	close(configuration.Done)
	*/
}

//
func loadSessionConfiguration(configurationFileName string) (configuration*SessionConfiguration) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
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

	return configuration
}

func setLogOutput() {
	f, err := os.Create("log.txt")
	if err != nil {
		fmt.Println("issue %s", err)
		return
	}
	log.SetOutput(f)
}

func SerialiseStruct(v interface{}) {
	defer func() {
		//swallow issue anc carry on
		if err := recover(); err != nil {
			fmt.Printf("local deffer %+v", err)
		}
	}()
	bytes, _ := json.Marshal(v)
	//owner=read+wqrite, group and others=read
	ioutil.WriteFile(DefaultStatusFileName, bytes, 0644)
}
func Printf(format string,args... interface{}){
	const (MainLiteral ="MAIN:")
	if len(args)<0 {
		fmt.Printf(MainLiteral+format+"\n")
	}else{
		fmt.Printf(MainLiteral+format+"\n",args)
	}
}

func usage(){
	fmt.Println("\nData integration tool. Usage:")
	fmt.Println("pupdate -config=monitoring/configuration_file [-test_config]")
	fmt.Println("-config path to configiration file name")
	fmt.Println("-test_config whether to test configuration")
	os.Exit(1)
}