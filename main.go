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
	//is there is any task to do
	if len(configuration.Tasks)==0{
		fmt.Println("no tasks scheduled. exit")
		os.Exit(1)
	}

	fmt.Printf("Configuration loaded\n")
	var jobDataChannel = make(chan *JobData, configuration.Tasks[0].Concurrency)
	var JobController = JobController{JobDataChannel: jobDataChannel}
	var taskId uint64 = 0

	JobController.StartTask(taskId, &configuration.Tasks[0])

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

//tested on MySQL updated
func (this *JobController) SQLUpdate(jobData *JobData, dsn string, sessionParams string, query string) {
	//store start time
	jobData.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			//log.Println(err)
			jobData.Error = true
			//jobData.ErrorMsg=err.Error()
		}
		//increase number of attmpts
		jobData.Attempts++
		//record data
		jobData.StopTime = time.Now()
		//notify producer that another job has finished
		this.JobDataChannel <- jobData
	}()

	//fmt.Printf("start: job_id=%d, start_id=%d stop_id=%d\n",jobId, startid, limit)

	//how to use connection pool?
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//log.Print("connection open ", Dsn)

	if len(sessionParams) > 0 {
		_, err = db.Exec(sessionParams)
	}

	if err != nil {
		panic(err)
	}

	//all data source details should be well encapsulated
	_, err = db.Exec(query)

	//log.Print("query finished:", query)

	if err != nil {
		panic(err)
	}

	//fmt.Printf("end: job_id=%d,start_id=%d stop_id=%d\n",jobId, startid, limit)
}
