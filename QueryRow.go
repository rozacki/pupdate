package main

import (
	"database/sql"
)


///Atomic SQL exec.
//tested on MySQL update
//todo:handle errors-
//todo:implement query
//todo: implement queryrow - that returns at most one row
//Executes a query and returns single value
func (this *SessionController) QueryRow(jobContext JobExecutionContextInterface){
	const MethodName	=	"QueryRow"

	//log some debuf info if in debug mode
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
			this.Debugf(jobContext.IsDebug(),"panic %s\n",err.Error())
			jobContext.SetErrorMessage(err.Error())
		}
		//increase number of attmpts
		jobContext.IncreaseAttempts()
		//record data
		jobContext.StopTime()
		//notify producer that another job has finished
		jobContext.Finish()
		//
		this.Debugf(jobContext.IsDebug(),FormatStruct,jobContext)
		this.Debugf(jobContext.IsDebug(),FormatStruct,jobContext.GetJobData())
	}()

	//how to use connection pool?
	db, err := sql.Open("mysql", jobContext.GetDsn())
	if err != nil {
		panic(err)
	}
	//close connection hence no side effects
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
	//result, err = db.Exec(jobContext.JobData.Query)
	var value interface{}
	err =	db.QueryRow(jobContext.GetQuery()).Scan(value)
	jobContext.SetValue(value)

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