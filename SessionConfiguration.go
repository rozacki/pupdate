package main
import "fmt"

type SessionConfiguration struct{
	//monitoring
	Monitoring MonitoringConfiguration
	//notificactions
	Notifications NotificationConfiguration
	//all tasks for this session
	Tasks []TaskConfiguration
	TaskCounter	uint64
	//
	Done chan struct{}
}

func (this *SessionConfiguration) SessionSuccess()(error){
	return GMonitoring.TraceOK("",SessionSuccess,true)

}
func (this *SessionConfiguration) SessionFailed(reason string)(error){
	return GMonitoring.Tracef("%s, reason:%s",SessionFailed,reason)
}

func (this *SessionConfiguration) TaskDisabled(taskName string)(error){
	return GMonitoring.Trace(taskName,TaskDisabled)
}

func (this*SessionConfiguration) Init()error{
	for i,_:=range this.Tasks{
		if err:=this.Tasks[i].Init();err!=nil{
			return err
		}
	}
	return nil;
}