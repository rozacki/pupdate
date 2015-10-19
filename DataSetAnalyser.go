package main

//provides, min, max and count of rows in database
//also provides types and distribution to help build more efficient queries
type DataSetAnalyser	interface{
	//
	Analyse(dsn string) (error,DataSetInfo)
	//
	Compare(dsi1, dsi2 DataSetInfo) (error,bool)
}

type dataSetAnalyser struct{

}

func (this*dataSetAnalyser)Analyse(dsn string) (dsi *DataSetInfo,err error){

	dsi.Name	=	dsn
	dsi.Source	=	dsn
	dsi.Size	=	0

/*
	Hash string		``
	//

	//
	Size	int64	``
	//
	Schema	string	``
	//
	Heat	string	`help: "how id is distributed in this data set.`
	//
	Parent	string	`help: "previous version"`
	//
	Max		int64
	//
	Min		int64
*/
	return nil,nil
}
//
func (this*dataSetAnalyser) Compare(dsi1, dsi2 DataSetInfo) (b bool,err error){
	return false,nil
}

