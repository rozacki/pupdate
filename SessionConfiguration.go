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
	return Monitoring.Event(this.SessionID,"",0,0,StartSession,this)
}

func (this *SessionConfiguration) StopSession()(*MonitoringError){
	return Monitoring.Event(this.SessionID,"",0,0,StopSession,this)
}