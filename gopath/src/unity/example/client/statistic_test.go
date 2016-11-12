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
	var st GoStatistic
	for i := 0; i < 100; i++ {
		if !GetStat(&st) {
			t.Error("GetStat() returns unexpected false")
		}
		//fmt.Println(i, st.NumGC)
		time.Sleep(time.Millisecond * 10)
	}
	//StopRoom()
	fmt.Println(st.NumGC)
	StopStatistic()
	if GetStat(&st) {
		t.Error("GetStat() returns unexpected true")
	}
}
