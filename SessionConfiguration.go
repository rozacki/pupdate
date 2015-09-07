package main

type SessionConfiguration struct{
	//monitoring
	SessionMonitoring Monitoring
	//notificactions
	SessionNotifications Notifications
	//all tasks for this session
	Tasks []TaskConfiguration
}
