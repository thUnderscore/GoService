package client

//#include "dataobj.h"
import "C"

import (
	"runtime"
	"time"
	"unity/example/shared"
)

var l GoStatistic
var n GoStatistic
var b GoStatistic
var last = &l
var next = &n
var buf = &b

//GoStatistic  container for go statistic
type GoStatistic C.struct_GoStatisticTag

type startCommand struct {
	d    time.Duration
	resp chan (byte)
}

var start = make(chan *startCommand)
var stop = make(chan chan byte)
var get = make(chan byte)
var ret = make(chan *GoStatistic)
var memStats = &runtime.MemStats{}

//CollectStatistic Collects statistic and populates structre's fields
func (st *GoStatistic) CollectStatistic() {
	if st == nil {
		return
	}

	runtime.ReadMemStats(memStats)

	st.NumGoroutine = C.int(runtime.NumGoroutine())
	st.Alloc = C.uint64_t(memStats.Alloc)
	st.Mallocs = C.uint64_t(memStats.Mallocs)
	st.Frees = C.uint64_t(memStats.Frees)
	st.HeapAlloc = C.uint64_t(memStats.HeapAlloc)
	st.StackInuse = C.uint64_t(memStats.StackInuse)
	st.PauseTotalNs = C.uint64_t(memStats.PauseTotalNs)
	st.NumGC = C.uint64_t(memStats.NumGC)

}

var gcc int

func statisticManager() {
	var cntr int
	stopCol := make(chan byte)
	inCol := make(chan *GoStatistic)
	outCol := make(chan *GoStatistic)

	for {
		select {
		case cmd := <-start:
			cntr++
			if cntr == 1 {
				shared.Logf("Start statistic collection. Interval: %v", cmd.d)
				last.Interval = C.int64_t(cmd.d)
				next.Interval = last.Interval
				buf.Interval = last.Interval

				last.CollectStatistic()
				go statisticCollector(stopCol, cmd.d, next, inCol, outCol)
			}
			cmd.resp <- 1
		case resp := <-stop:
			if cntr > 0 {
				cntr--
				if cntr == 0 {
					shared.Log("Stop statistic collection")
					stopCol <- 1
				}
			}
			resp <- 1
		case <-get:
			if cntr == 0 {
				ret <- nil
			}
			gcc++

			last.InUse = 1
			buf.InUse = 0
			ret <- last

		case st := <-outCol:
			if last.InUse == 1 {
				tmp := buf
				buf = last
				last = st
				inCol <- tmp
			} else {
				tmp := last
				last = st
				inCol <- tmp
			}

			if gcc == 10 {
				gcc = 0
				runtime.GC()
			}

		}
	}
}

func statisticCollector(st chan byte, d time.Duration, m *GoStatistic, in chan *GoStatistic, out chan *GoStatistic) {
	ticker := time.NewTicker(d)
	c := ticker.C
	n := m
	for {
		select {
		case <-st:
			return
		case <-c:
			n.CollectStatistic()
			out <- n
			n = <-in
			//Log("COLLECTED")
		}
	}
}

//StartStatistic starts collection of statistic with given interval (in ms)
//export StartStatistic
func StartStatistic(interval int) {
	cmd := startCommand{d: time.Millisecond * time.Duration(interval), resp: make(chan byte)}
	start <- &cmd
	<-cmd.resp
}

//StopStatistic stops collection ofstatistic
//export StopStatistic
func StopStatistic() {
	resp := make(chan byte)
	stop <- resp
	<-resp
}

//GetStat returns pointer to last collected statistic. Expected to becalled in a loop by ONLY ONE consumer
//export GetStat
func GetStat() *GoStatistic {

	get <- 1
	return <-ret
}
