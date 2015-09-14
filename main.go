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

//will be used to notify users about some important facts
var Notifier *NotificationsModule
//global module for recording session,task, job progress
var Monitoring Monitor
//last sucessfull etl
//todo: change name to something more generic. move it to session context when available
var LastEtl	time.Time

func main() {
	//parse parameter
	var configFileName = flag.String("config", "", "configuration file")
	flag.Parse()

	if len(*configFileName) == 0 {
		usage()
	}

	//if file is corrupted then we get nil
	//if types dont match then we get zero value for that SessionConguration
	configuration := loadSessionConfiguration(*configFileName)
	if configuration == nil {
		fmt.Println("configuration load error")
		os.Exit(1)
	}
	Printf("Configuration loaded. Found %d tasks\n", len(configuration.Tasks))

	if Monitoring=makeMonitoring(configuration.Monitoring);Monitoring==nil{
		os.Exit(1)
	}

	//find last_elt timestamp
	if lastEtl,err:=findLastEtlTime(); err!=nil{
		Printf("error finding $LastEtl:%s\n", err)
		os.Exit(1)
	}else{
		LastEtl=lastEtl
		Printf("$LastEtl = %s\n",LastEtl.Format(LastEtlFileFormat))
	}

	configuration.SessionID	=	time.Now().Format(SessionFileFormat)
	configuration.Done		=	make(chan struct{})
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
	fmt.Printf("main:"+format+"\n",args)
}

func usage(){
	fmt.Println("Data integration tool")
	fmt.Println("Usage:")
	fmt.Println("pupdate -config=monitoring/configuration_file")
	os.Exit(1)
}