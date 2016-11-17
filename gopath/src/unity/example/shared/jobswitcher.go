package shared

import "sync"
import "fmt"

//Job is interface that describe wrapper for goroutine
type Job interface {
	Start(once bool)
	Stop(sync bool)
	IsActive() bool
}

//JobSwitcher extended switcher that implements Job interface and
type JobSwitcher struct {
	*Switcher
	//Cond used for goroutine loop  implementation
	*sync.Cond
	syncStop bool
	OnStart  func()
	OnStop   func()
}

//Start starts job. You can't start job if it's already started and not stoped
//If job was stoped you should pass once == false if you want restart job
func (s *JobSwitcher) Start(once bool) {
	s.On(s.OnStart, once)
}

//Stop stops job
func (s *JobSwitcher) Stop(sync bool) {
	s.Off(func() {
		s.syncStop = sync
		if s.OnStop != nil {
			s.OnStop()
		}
		if s.syncStop {
			s.Wait()
			if s.Active {
				fmt.Println("WTF")
			}
		}
	})
}

//RunGoroutine starts f in goroutine, waits for actual start
//goroutine signal if job was stopped in sync mode
//if ownLockInF == true, goroutine owns lock before f call and after f finished
func (s *JobSwitcher) RunGoroutine(f func(), ownLockInF bool) {
	go func() {
		//signal InternalStart
		s.Lock()
		s.Signal()
		if !ownLockInF {
			s.Unlock()
		}
		f()
		//signal to InternalStop
		if !ownLockInF {
			s.Lock()
		}
		//just in case
		s.Active = false
		if s.syncStop {
			s.Signal()
		}
		s.Unlock()
	}()
	s.Wait()
}

//NewJobSwitcher creates and initializes JobSwitcher
func NewJobSwitcher() *JobSwitcher {
	sw := &JobSwitcher{Switcher: NewSwitcher()}
	sw.Cond = sync.NewCond(sw.Mutex)
	return sw
}
