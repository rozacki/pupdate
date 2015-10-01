package main

type SessionConfiguration struct{
	//monitoring
	Logging MonitoringConfiguration
	//notificactions
	Notifications NotificationConfiguration
	//all tasks for this session
	Tasks []TaskConfiguration
	//
	TaskCounter	uint64
	//interface to session execution context
	executionContext SessionExecutionContext
	//
	Test interface{}
}

func (this*SessionConfiguration) Init()error{
	for i,_:=range this.Tasks{
		if err:=this.Tasks[i].Init();err!=nil{
			return err
		}
	}
	return nil;
}
//
func (this *SessionConfiguration) SessionSuccess()(error){
	return GLogger.Tracef(SessionSuccess)

}
func (this *SessionConfiguration) SessionFailed(reason string)(error){
	return GLogger.Tracef("%s, reason:%s",SessionFailed,reason)
}

func (this *SessionConfiguration) TaskDisabled(taskName string)(error){
	return GLogger.Tracef("%s %s", taskName,TaskDisabled)
}