package client

//#include "dataobj.h"
import "C"
import (
	"runtime"
	"time"
	"unity/example/shared"
)

//GoStatistic  container for go statistic
type GoStatistic C.struct_GoStatisticTag

//StatisticMan implements statistic collector. Expected to be used by ONLY ONE consumer
type StatisticMan struct {
	*shared.JobChan
	memStats *runtime.MemStats
	last     *GoStatistic
	next     *GoStatistic
	interval time.Duration
}

//Get populate passed struct by last collected statistic
func (s *StatisticMan) Get(res *GoStatistic) bool {
	s.Lock()
	if !s.Active {
		s.Unlock()
		return false
	}
	*res = *s.last
	s.Unlock()
	return true

}

//collect Collects statistic and populates structre's fields
func (s *StatisticMan) collect(st *GoStatistic) {
	runtime.ReadMemStats(s.memStats)
	st.NumGoroutine = C.int(runtime.NumGoroutine())
	st.Alloc = C.uint64_t(s.memStats.Alloc)
	st.Mallocs = C.uint64_t(s.memStats.Mallocs)
	st.Frees = C.uint64_t(s.memStats.Frees)
	st.HeapAlloc = C.uint64_t(s.memStats.HeapAlloc)
	st.StackInuse = C.uint64_t(s.memStats.StackInuse)
	st.PauseTotalNs = C.uint64_t(s.memStats.PauseTotalNs)
	st.NumGC = C.uint64_t(s.memStats.NumGC)
}

//NewStatisticMan create new instance of StatisticMan
func NewStatisticMan() *StatisticMan {
	s := &StatisticMan{
		next:     new(GoStatistic),
		last:     new(GoStatistic),
		memStats: &runtime.MemStats{},
	}

	s.JobChan = shared.NewJobChan(func(j *shared.JobChan) {
		ticker := time.NewTicker(s.interval)
		c := ticker.C
		intv := C.int64_t(s.interval)
		s.last.Interval = intv
		s.next.Interval = intv
		s.collect(s.last)
		for {
			select {
			case <-j.ExitChn:
				ticker.Stop()
				return
			case <-c:
				s.collect(s.next)
				j.Lock()
				s.next, s.last = s.last, s.next
				j.Unlock()
			}
		}
	})
	return s
}

//StartStatistic starts collection of statistic with given interval (in ms)
//export StartStatistic
func StartStatistic(interval int) {
	if stat == nil {
		stat = NewStatisticMan()
	} else {
		if stat.IsActive() {
			return
		}
	}
	stat.interval = time.Duration(interval) * time.Millisecond
	stat.Start()

}

//StopStatistic stops collection ofstatistic
//export StopStatistic
func StopStatistic() {
	if stat == nil {
		return
	}
	stat.Stop(true)
}

//GetStat populate passed struct by last collected statistic
//export GetStat
func GetStat(res *GoStatistic) bool {
	return stat.Get(res)
}
