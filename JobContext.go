package main

type JobContext struct{
	JobData *JobData
	Dsn string
	SessionParams string
	JobDataChannel chan *JobData
	Debug bool
}
