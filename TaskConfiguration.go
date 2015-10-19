package main

import(
	"strings"
	"reflect"
)

const(
	Exec				=	"EXEC"
	QueryOne			=	"QueryOne"
	Query				=	"Query"
	HelpTag				=	"help"
	HelpExampleTag		=	"example"
	JsonValueSeparator 	=	"\n"
)

type TaskConfiguration struct {
	Name   string			`help:""`
	Dsn    string			`help:""`
	Exec   string			`help:"SQL query to execute"`
	ExecTab[] string 		`help:"Tab version may be used when we don't want to type one query as long line."`
	ExecType string			`help:"Tap of query to execute:queryone - scalar result, query	-vector result,exec - no result, default"`
	//
	Output string				`help:"file_name|stdout, stdout is default For queryone and query we can store output into files or memory..."`
	//
	MaxQuery	string				`help:"if type string it may be sql query, if int then it is max value"`
	Max 	uint64					`help:""`
	//
	MinQuery	string				`help:"If type string it may be sql query, if int then it is min value"`
	Min         uint64				`help:""`
	Step        uint64 				`help:"Used during data transferring base don some id`
	//
	Concurrency uint64			`help:"how many concurrent sql connections to open"`
	//
	PreSteps []string			`help:"one of many possible parameters of session"`
	//
	SessionParamTab []string	`help:"one of many possible parameters of session"`
	//
	Method string				`help:"name of method to run. currenlty supported: SQLUpdate,TestSQL. Default is SQLUpdate"`
	//
	MaxAttempts uint64			`help:"How many attempts before job is dropped"`
	//
	Disabled bool				`help:""`
	//
	Debug bool					`help:"Allows to debug individual tasks. This flag is inherited by JonContext."`
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
func (this TaskConfiguration) GetHelpString(format string) string{
	helpString:=""
	if format==""{
		format=JsonValueSeparator
	}
	tt:=reflect.TypeOf(this)
	for i:=0;i< tt.NumField();i++{
		if helpString!=""{
			helpString+=format
		}
		helpString+=tt.Field(i).Name
		helpString+=":\n\t"
		helpString+=tt.Field(i).Tag.Get(HelpTag)
	}
	return helpString
}