package client

//#include "dataobj.h"
import "C"

import (
	"runtime"
	"time"
	"unity/example/shared"
)

type Statistic struct {
	*shared.JobChan
	memStats *runtime.MemStats
	last     *GoStatistic
	next     *GoStatistic
	buf      *GoStatistic
	Interval time.Duration
}

//Get returns pointer to last collected statistic. Expected to be called in a loop by ONLY ONE consumer
func (s *Statistic) Get() *GoStatistic {
	s.Lock()
	if !s.Active {
		s.Unlock()
		return nil
	}
	res := s.last
	s.last.InUse = 1
	s.buf.InUse = 0
	s.Unlock()
	return res
}

//collect Collects statistic and populates structre's fields
func (s *Statistic) collect(st *GoStatistic) {
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

func NewStatistic() *Statistic {
	s := &Statistic{
		next:     new(GoStatistic),
		last:     new(GoStatistic),
		buf:      new(GoStatistic),
		memStats: &runtime.MemStats{},
	}

	s.JobChan = shared.NewJobChan(func(j *shared.JobChan) {
		ticker := time.NewTicker(s.Interval)
		c := ticker.C
		intv := C.int64_t(s.Interval)
		s.last.Interval = intv
		s.next.Interval = intv
		s.buf.Interval = intv
		for {
			select {
			case <-j.ExitChn:
				ticker.Stop()
				return
			case <-c:
				s.collect(s.next)
				j.Lock()
				if s.last.InUse == 1 {
					s.last.InUse = 0
					s.next, s.buf, s.last = s.buf, s.last, s.next
				} else {
					s.next, s.last = s.last, s.next
				}
				j.Unlock()
				//Log("COLLECTED")
			}
		}
	})
	return s
}
