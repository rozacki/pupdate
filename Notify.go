package main

//notification module: will use whatever means to notify user about the result
//todo: change name to notifier
//todo: use standard syslogger’
//todo:provide reference to local log files in messages
type NotificationsModule struct{
	Configuration NotificationConfiguration
}

type NotificationConfiguration struct{
	Email string
}