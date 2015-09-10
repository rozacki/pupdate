package main

type JobContext struct{
	JobData *JobData
	Dsn string
	SessionParams string
	Query string
	JobDataChannel chan *JobData
}
