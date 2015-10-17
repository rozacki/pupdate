package main

import (
	"time"
)


//Defines single unit of job
type JobData struct {
	Name      string
	Id        uint64
	StartTime time.Time
	StopTime  time.Time
	Error     bool
	LastErrorMsg string
	Attempts  uint64
	//
	Query string
	//how many rows have been affected based on what driver returns
	RowsAffected uint64
}
