package client

//#include "dataobj.h"
import "C"

import "unity/example/shared"

//"context"
//"google.golang.org/grpc"

//Connector  representsconnection point for client
type Connector C.struct_ClientConnectorTag

//Conn represent connetion point between go lib and client
var Conn *Connector

//Svcs represent root group of services
var Svcs *shared.ServiceGroup

//StartClient starts client. Call it at the very beginning of app life cycle. Not thread safe
//export StartClient
func StartClient(c *Connector) {
	if Conn != nil {
		shared.Log("Client is already started")
	}
	Conn = (*Connector)(c)
	Svcs = shared.NewServiceGroup()
	shared.SetLogger(&clientLogger{c.log})
	shared.Log("StartClient")
	go statisticManager()

}

//StopClient stops client. Call it at the very end of app life cycle. Not thread safe
//export StopClient
func StopClient() {
	if Conn == nil {
		shared.Log("Client is not started")
		return
	}
	shared.Log("StopClient")
	Svcs.StopGroup(nil)
	shared.ResetLogger()
	Svcs = nil
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

/*
	const (
		address     = "localhost:50051"
		defaultName = "world"
	)

		gconn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			Log(fmt.Sprintf("did not connect: %v", err))
		}
		defer func() {
			defer gconn.Close()
			if r := recover(); r != nil {
				Log(fmt.Sprintf("Recovered: %v", r))
			}
		}()

		gc := NewGreeterClient(gconn)

		// Contact the server and print out its response.
		name := defaultName

		r, err := gc.SayHello(context.Background(), &HelloRequest{Name: name})
		if err != nil {
			Log(fmt.Sprintf("could not greet: %v", err))
		} else {
			Log(fmt.Sprintf("Greeting: %s", r.Message))
		}
*/
