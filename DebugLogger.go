package main

import "fmt"

//specialised logger for debugging session
//If we pass all calls through this interface and implementation then it will be easier to change in future instead of just using fmt...
type DebugLogger interface{
	//this is generic interface that is not called direclty but by wrapping functions that output context specific string for Session, Tasks, Jobs
	Debugf(format string, args []interface{})
}

type DebugLoggingModule struct{
}

func (this*DebugLoggingModule) Debugf(format string, args []interface{}){
	const (DebugLiteral="DEBUG:")
	defer func(){
		//just dummy recover and carry on
		recover()
	}()

	s:=DebugLiteral+format
	if len(args)==0{
		fmt.Println(s)
	}else{
		fmt.Printf(s+"\n",args...)
	}
}