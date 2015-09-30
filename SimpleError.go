package main

//Error information specific to Monitoring module
type LoggingError struct {
//original error
error
Msg      string
}

func (this LoggingError) Error() string{
return this.Msg
}