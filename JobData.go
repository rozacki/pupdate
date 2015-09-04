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
	ErrorMsg  string
	Attempts  uint64
	//parrtition begining
	PartStart uint64
	//partition end
	PartEnd uint64
}
