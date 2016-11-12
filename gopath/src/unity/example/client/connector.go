package client

//#include "dataobj.h"
import "C"

import "unity/example/shared"

//Connector  representsconnection point for client
type Connector C.struct_ClientConnectorTag

//Conn represent connetion point between go lib and client
var Conn *Connector
var stat *StatisticMan

//StartClient starts client. Call it at the very beginning of app life cycle. Not thread safe
//export StartClient
func StartClient(c *Connector) {
	if Conn != nil {
		shared.Log("Client is already started")
	}
	Conn = (*Connector)(c)
	shared.SetLogger(&clientLogger{c.log})
	shared.Log("StartClient")

}

//StopClient stops client. Call it at the very end of app life cycle. Not thread safe
//export StopClient
func StopClient() {
	if Conn == nil {
		shared.Log("Client is not started")
		return
	}
	shared.Log("StopClient")
	shared.ResetLogger()
	Conn = nil
}

//StartRoom prepare
//export StartRoom
func StartRoom() {

}

//StopRoom finalize room
//export StopRoom
func StopRoom() {

}

//Count increase count and returns result
//export Count
func Count() int {

	if Conn == nil {
		return -1
	}
	Conn.counterValue = Conn.counterValue + 1

	res := Conn.counterValue
	/*
		if res > 10300 {
			conn = nil
		}
	*/
	//conn.Log("Log from GO")
	return int(res)

}
