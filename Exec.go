package main

import (
	"time"
	"database/sql"
)


///Atomic SQL exec.
//tested on MySQL update
//todo:handle errors-
//todo:implement query
//todo: implement queryrow - that returns at most one row
//Exec executes a query without returning any rows
func (this *SessionController) Exec(jobContext JobExecutionContext) (err error){
	//log some debuf info if in debug mode
	this.Debugf(jobContext.Debug,FormatStruct,jobContext)
	this.Debugf(jobContext.Debug,FormatStruct,jobContext.JobData)

	//store start time- for reporting purposes
	jobContext.JobData.StartTime = time.Now()

	//all defered functions
	defer func() {
		if recovered := recover(); err != nil {
			switch v:=recovered.(type){
				case error: err=v
				default: err=StringError("unknown error")
			}
			this.Debugf(jobContext.Debug,"panic %s\n",err.Error())
			jobContext.JobData.Error 		= 	true
			jobContext.JobData.LastErrorMsg	=	err.Error()
		}
		//increase number of attmpts
		jobContext.JobData.Attempts++
		//record data
		jobContext.JobData.StopTime = time.Now()
		//notify producer that another job has finished
		jobContext.JobDataChannel <- jobContext.JobData
		//
		this.Debugf(jobContext.Debug,FormatStruct,jobContext)
		this.Debugf(jobContext.Debug,FormatStruct,jobContext.JobData)
	}()

	//how to use connection pool?
	db, err := sql.Open("mysql", jobContext.Dsn)
	if err != nil {
		panic(err)
	}
	//close connection hence no side effects
	defer db.Close()
	//log.Print("connection open ", Dsn)

	//iterate all 'set'
	for _,stmt:=range jobContext.PreSteps {
		_, err := db.Exec(stmt)
		if err != nil {
			panic(err)
		}
		this.Debugf(jobContext.Debug,"pre-exec: %s",stmt)
	}

	var result sql.Result
	//all data source details should be well encapsulated
	result, err = db.Exec(jobContext.JobData.Query)

	if err != nil {
		panic(err)
	}
	//if driver supports rows affected and last inserted id
	if rowsAffected,err:=result.RowsAffected();err==nil{
		jobContext.JobData.RowsAffected=uint64(rowsAffected)
	}

	this.Debugf(jobContext.Debug,"exec query %s \n",jobContext.JobData.Query)
	this.Debugf(jobContext.Debug,"rows affected: %d\n",jobContext.JobData.RowsAffected)

	return nil
}