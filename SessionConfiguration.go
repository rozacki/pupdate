package main

type SessionConfiguration struct{
	//monitoring
	Monitoring MonitoringConfiguration
	//notificactions
	Notifications NotificationConfiguration
	//all tasks for this session
	Tasks []TaskConfiguration
	//
	SessionID	string
	TaskCounter	uint64
	//
	Done chan struct{}
}

func (this *SessionConfiguration) StartSession()(*MonitoringError){
	return Monitoring.Trace(this.SessionID,"",0,0,StartSession,this)
}

func (this *SessionConfiguration) SessionSuccess()(*MonitoringError){
	return Monitoring.TraceOK(this.SessionID,"",0,0,StopSession,this,true)

}
func (this *SessionConfiguration) SessionFail()(*MonitoringError){
	return Monitoring.Trace(this.SessionID,"",0,0,StopSession,this)
}

func (this *SessionConfiguration) Trace(msg string)(*MonitoringError){
	return Monitoring.Trace(this.SessionID,"",0,0,Trace,msg)
}