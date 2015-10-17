package main

import(

	"strings"
)

const(
	Exec		=	"EXEC"
	QueryOne	=	"QueryOne"
	Query		=	"Query"
)

type TaskConfiguration struct {
	Name   string
	Dsn    string
	//single query no result
	Exec   string
	ExecTab[] string
	//queryone - scalar result
	//query	-vector result
	//exec - no result is default
	ExecType string
	//single query, vector result
	//for queryone and query we can store output into files or memory...
	Output string
	//if type string it may be sql query, if int then it is max value
	MaxQuery	string
	Max 	uint64

	//If type string it may be sql query, if int then it is min value
	MinQuery	string
	Min         uint64
	//
	Step        uint64
	Concurrency uint64
	//one of many possible parameters of session
	PreSteps []string
	//one of many possible parameters of session
	SessionParamTab []string
	//name of method to run. currenlty supported: SQLUpdate,TestSQL. Default is SQLUpdate
	Method string
	//How many attempts before job is dropped
	MaxAttempts uint64
	//
	Disabled bool
	//Allows to debug individual tasks. This flag is inherited by JonContext.
	Debug bool
//**************** dynamic parameters
// todo: should go to taskcontext
	//current task Id
	TaskId uint64
}
//TaskConfiguration extands existing MonitoringModule
func (this* TaskConfiguration) EventStartTask()(error){
	return GLogger.Tracef("%s %s",this.Name,StartTask)
}

func (this* TaskConfiguration) EventSuccessTask()(error){
	return GLogger.Tracef("%s %s",this.Name,TaskSuccess)
}

func (this* TaskConfiguration) EventFailTask(reason string)(error){
	return GLogger.Tracef("%s %s,reason: %s ",this.Name,TaskFailed,reason)
}

func (this* TaskConfiguration) Trace(data interface{})(error){
	return GLogger.Tracef("%s %s",this.Name,Trace)
}

func (this* TaskConfiguration) RowsAffected(rowsAffected uint64)(error){
	return GLogger.Tracef("%s %s:%d",this.Name,RowsAffected,rowsAffected)
}

func (this* TaskConfiguration) TotalRowsAffected(rowsAffected uint64)(error){
	return GLogger.Tracef("%s %s:%d",this.Name,TotalRowsAffected,rowsAffected)
}

//does some housekeeping int he task configuration, shoudl be called after configuration is loaded
func (this* TaskConfiguration) Init() error{
	if len(this.ExecTab)>0{
		//join ExecTab into Exec if Exec is emopty and ExecTab is not
		this.Exec=strings.Join(this.ExecTab,"")
		//fmt.Println(this.Exec)
	}
	return nil
}