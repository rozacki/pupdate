package main

func PrependArrayOfInterfaces(a []interface{},b []interface{})(c []interface{}){
	c=b[0:len(b)]
	for _,v:=range a{
		c=append(c,v)
	}
	return c
}
