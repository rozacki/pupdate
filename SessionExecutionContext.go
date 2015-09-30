package main

type SessionExecutionContext interface{
	//
	SessionSuccess()(error)
	SessionFailed(reason string)(error)
	TaskDisabled(taskName string)(error)
}


