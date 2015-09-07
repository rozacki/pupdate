package main


//comfoguration todo:
// 	multiple tasks from one configuration
//	add storage_type={mysql|postgress|mongodb...}
//	add type of operation:{update|insert...}
//todo error handling:
//	MySQL schema issues

import (
	"database/sql"
	_ "database/sql/driver"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"time"
	_ "time"
)

type JobConfiguration struct {
}

const (
	InvalidJob            = 0
	TaskMethod_SQLUpdate  = "SQLUpdate"
	TaskMethod_TestUpdate = "TestUpdate"
	DefaultStatusFileName = "tasks.json"
)

func main() {

	var configFileName = flag.String("config", "", "configuration file")
	flag.Parse()

	if len(*configFileName) == 0 {
		fmt.Print(*configFileName)
		fmt.Println("missing configuration (-config parameter) file path")
		os.Exit(1)
	}

	configuration := loadConfiguration(*configFileName)
	if configuration == nil {
		fmt.Println("configuration load error")
		os.Exit(1)
	}


	fmt.Printf("Configuration loaded. Found %d tasks", len(configuration.Tasks))

	var sessionController 	SessionController
	//start new session
	if sessionController:=	makeSession(configuration);sessionController==nil{
		os.Exit(1)
	}
	sessionController.StartTask(&configuration.Tasks[0])

	var key string
	fmt.Scanf("%s", &key)
}

//
func loadConfiguration(configurationFileName string) (configuration*SessionConfiguration) {
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

