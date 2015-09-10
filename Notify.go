package main

//notification module: will use whatever means to notify user about the result
//todo: change name to notifier
type NotificationsModule struct{
	Configuration NotificationConfiguration
}

type NotificationConfiguration struct{
	Email string
}