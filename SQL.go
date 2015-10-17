package main
import "database/sql"

type SQLJobber interface{
	StartJob(dsn string) (error)
	StopJob()
	//if sql driver supports it will return rowsaffected and last id
	Exec(query string, params... interface{}) (rowsAffected int, lastId int,err error)
	//is it complete?
	QueryRow(query string, params... interface{}) (interface{},error)
	//this method proto is not complete
	Query(query string, params... interface{}) (error)
	//return last error
	GetLastError() error
	//is it working
	Status() bool
}

type SQLJob struct {
	dsn string
	query string
	params []interface{}
	working bool
	lastError error

	db*sql.DB

}
//Opens a new connection- it actually opens a new connection pool
func (this*SQLJob) StartJob(dsn string)(err error){
	if this.db,err=sql.Open("mysql",dsn);err!=nil{
		//log and exit
		return err
	}
	return nil
}
