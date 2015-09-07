package main

import (
	"fmt"
	"math/rand"
	_ "os"
	"path"
	"path/filepath"
	"testing"
	"time"
)

const (
	TestConlfigurationFileName = "test_configs/test1.json"
	TestConfigurationFolder    = "test_configs"
)

//todo: load and run all configurations from /test folder
func TestStartTask(T *testing.T) {

	configs, err := filepath.Glob(path.Join(TestConfigurationFolder, "*"))
	if err != nil {
		T.Fatal(err.Error())
	}
	for _, configFileName := range configs {
		configuration := loadConfiguration(configFileName)
		if configuration == nil {
			T.Logf("configuration load error %s\n", configFileName)
			continue
		}
		if configuration.Disabled {
			continue
		}

		fmt.Printf("Configuration loaded: %s from file %s\n", configuration.Name,configFileName)
		var jobDataChannel = make(chan *JobData, configuration.Concurrency)
		var JobController = JobController{JobDataChannel: jobDataChannel}
		var taskId uint64 = 0

		//test in serial
		taskData := JobController.StartTask(taskId, configuration)
		//task finished, we can verify task data
		VerifyTestOutcome(T, taskData, *configuration)
	}
}
func VerifyTestOutcome(T *testing.T, taskData TaskData, configuration TaskConfiguration) {
	var failed bool
	//test outcome and raise a flag
	if taskData.QueueLength > 0 {
		failed=true
		T.Errorf("Test failed, reason queue not empty\n")
	}

	if taskData.Success+taskData.MaxAttemptJobsDropped != ((configuration.Max - configuration.Min)/configuration.Step) {
		T.Errorf("Test failed, reason not all items processed\n")
		failed=true
	}

	if taskData.Status != "finished" {
		T.Errorf("Test failed, reason job status different than 'finished\n")
		failed=true
	}
	if !failed{
		T.Logf("PASS")
	}
	T.Logf("TaskData: %+v", taskData)
}

//mock up simplest etst update
//succsefull and unsucsesfull jobs
//todo: time based
func (this *JobController) TestUpdate(jobData *JobData, dsn string, sessionParams string, query string) {
	//succsefull and unsacsesfull jobs
	jobData.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic")
			jobData.Error = true
		}
		//increase number of attmpts
		jobData.Attempts++
		//record data
		jobData.StopTime = time.Now()
		//notify producer that another job has finished
		this.JobDataChannel <- jobData
	}()
	//randomly panic, but use the same seed hence the same results
	if rand.Intn(100) < 90 {
		panic("every 10th attempt is error")
	}
}
