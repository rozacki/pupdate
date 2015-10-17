package main


//provides, min, max and count of rows in database
//also provides types and distribution to help build more efficient queries
type DataSetAnalyser	interface{
	GetMax() (uint64)
	GetMin() (uint64)
	GetType() (interface{})
}

type DataSetInfo	struct{
	Max	interface{}
	Min interface{}
	Dsn string
}
