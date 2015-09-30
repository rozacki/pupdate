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
	ExecutionContext SessionExecutionContext
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
	return GLogging.Tracef("",SessionSuccess)

}
func (this *SessionConfiguration) SessionFailed(reason string)(error){
	return GLogging.Tracef("%s, reason:%s",SessionFailed,reason)
}

func (this *SessionConfiguration) TaskDisabled(taskName string)(error){
	return GLogging.Tracef(taskName,TaskDisabled)
}