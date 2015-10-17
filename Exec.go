package main

import (
	"database/sql"
)

///Atomic SQL exec.
//tested on MySQL update
//todo:handle errors-
//todo:implement query
//todo: implement queryrow - that returns at most one row
//Exec executes a query without returning any rows
func (this *SessionController) Exec(jobContext JobExecutionContextInterface){
	//log some debug info if in debug mode
	this.Debugf(jobContext.IsDebug(),FormatStruct,jobContext)
	this.Debugf(jobContext.IsDebug(),FormatStruct,jobContext.GetJobData())

	//store start time- for reporting purposes
	jobContext.StartTime()

	//all defered functions
	defer func() {
		var err error
		if recovered := recover(); recovered != nil {
			switch v:=recovered.(type){
				case error: err=v
				default: err=StringError("unknown error")
			}
			this.Debugf(jobContext.IsDebug(),"paniced %s\n",err.Error())
			jobContext.SetErrorMessage(err.Error())
		}
		//increase number of attempts
		jobContext.IncreaseAttempts()
		//record data
		jobContext.StopTime()
		//notify controller that another job has finished
		jobContext.Finish()
		//
		this.Debugf(jobContext.IsDebug(),FormatStruct,jobContext)
		this.Debugf(jobContext.IsDebug(),FormatStruct,jobContext.GetJobData())
	}()

	//We open here connection poll and this is not recommended way of using it but it tested and reasonably fast
	//if is to be changed, how to implement stateful connections?
	db, err := sql.Open("mysql", jobContext.GetDsn())
	if err != nil {
		panic(err)
	}
	//close connection pool
	defer db.Close()
	//log.Print("connection open ", Dsn)

	//iterate all 'set'
	for _,stmt:=range jobContext.GetPreSteps() {
		_, err := db.Exec(stmt)
		if err != nil {
			panic(err)
		}
		this.Debugf(jobContext.IsDebug(),"pre-exec: %s",stmt)
	}

	var result sql.Result
	//all data source details should be well encapsulated
	result, err = db.Exec(jobContext.GetQuery())

	if err != nil {
		panic(err)
	}
	//if driver supports rows affected and last inserted id
	if rowsAffected,err:=result.RowsAffected();err==nil{
		jobContext.SetRowsAffected(uint64(rowsAffected))
	}

	this.Debugf(jobContext.IsDebug(),"exec query %s \n",jobContext.GetQuery())
	this.Debugf(jobContext.IsDebug(),"rows affected: %d\n",jobContext.GetRowsAffected())
}