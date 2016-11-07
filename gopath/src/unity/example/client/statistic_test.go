package client

import (
	"fmt"
	"testing"
	"time"
)

var c = Connector{}

func TestStatistic(t *testing.T) {
	StartClient(&c)

	//StartRoom(&c)
	StartStatistic(10)
	var st *GoStatistic
	for i := 0; i < 2000; i++ {
		st = GetStat()
		if st == nil {
			t.Error("GetStat() returns unexpected nil")
			return
		}
		//fmt.Println(i, st.NumGC)
		time.Sleep(time.Millisecond * 10)
	}
	//StopRoom()
	fmt.Println(st.NumGC)
	StopStatistic()
	st = GetStat()
	if st != nil {
		t.Error("GetStat() returns unexpected non nil value")
		return
	}
}
