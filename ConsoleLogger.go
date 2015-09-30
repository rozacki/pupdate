package main
import "fmt"

type ConsoleLogger interface{
	Printf(format string, args []interface{})
}

type ConsoleLoggerModule struct{

}

func (this*ConsoleLoggerModule)Printf(format string, args []interface{}){
	fmt.Printf(format,args)
}