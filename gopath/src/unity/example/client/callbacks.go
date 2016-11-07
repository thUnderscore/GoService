package client

//#include "dataobj.h"
import "C"

import "unsafe"

//CallStringCallback Calls callback if it's not empty and len(str) > 0
func CallStringCallback(callBack C.stringCallback, str string) {
	if callBack != nil && len(str) > 0 {
		bts := []byte(str)
		C.CallStringCallback(callBack, unsafe.Pointer(&bts[0]), C.int(len(str)))
	}
}
