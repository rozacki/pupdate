package main


//Meta-data about data set. It will be used to compare data sets
type DataSetInfo struct{
	//
	Name string		``
	//
	Hash string		``
	//
	Source string	``
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
}

