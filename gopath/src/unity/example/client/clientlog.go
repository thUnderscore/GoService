package client

//#include "dataobj.h"
import "C"

type clientLogger struct {
	cbck C.stringCallback
}

//Log send string to client's log
func (lgr *clientLogger) Log(str string) {
	CallStringCallback(lgr.cbck, str)
}
