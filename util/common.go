package util

import "log"

const (
	DEBUG = true
	//DEBUG = false
)

func DPrintln(args... interface{}){
	if DEBUG{
		log.Println(args...)
	}
}
func DPrintf(format string, args... interface{}){
	if DEBUG{
		log.Printf(format,args...)
	}
}

func Min(a,b int)int{
	if a<b{
		return a
	}
	return b
}

