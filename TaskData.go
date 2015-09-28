package main

import (
	"time"
)

type TaskData struct {
	Id          uint64
	Success     uint64
	Errors      uint64
	DataCursor uint64
	QueueLength uint64
	//current timeout used for waiting for job to report status
	Timeout int
	//copy of configuration
	TaskConfiguration TaskConfiguration
	//last time serialised
	Serialised time.Time
	// what is the curent status: {pending, sowkring,finished
	Status string
	//
	Name string
	//dropped jobs after reaching max attempts
	MaxAttemptJobsDropped uint64
	//
	JobId uint64
	//todo: not used
	//Can be used to register last sucessfull ETL
	CreationTime time.Time
	//last dropped job
	LastDroppedJob JobData
	//total rows affected by this task
	RowsAffected	uint64
	//indicated if task failed or not
	Failed bool
}
