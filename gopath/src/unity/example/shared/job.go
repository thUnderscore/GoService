package shared

import "sync"

//Job is interface that describe wrapper for goroutine
type Job interface {
	Start()
	Stop(sync bool)
	IsActive() bool
}

type jobSwitcher struct {
	Switcher
	//Cond used for goroutine loop  implementation
	*sync.Cond
	syncStop bool
}

//JobCond is wrapper for goroutine that waits sync.Cond in loop
//goroutine should return if JobCond.isOn == false after Cond.Wait
//Cond could be signaled by Stop or by JobCond.Wake
type JobCond struct {
	//thread safe Job state. Switcher.Mutex used in JobCond.Wake
	jobSwitcher
	//if isEx == true JobCond.f should implement for loop like in JobCond.InternalStart
	//implement this loop yourself in f function if isEx == true
	//			var isWaiting bool
	//			for {
	//				if j.cntr == 0 && j.Active {
	//					//wait wake
	//					j.Wait()
	//				}
	//				isWaiting = j.isWaiting
	//				j.isWaiting = false
	//				//if Stop and no any Wake happend after last handling j.cntr is 0
	//				//handle situation when j.cntr > 0 and j.Active == false in f if you need
	//				if j.cntr > 0 {
	//					//reset state
	//					j.cntr = 0
	//					//handle wakes
	//					...
	//				}
	//
	//				if isWaiting {
	//					//if WakeSync was called and some caller wait for back signal
	//					j.waitCnd.L.Lock()
	//					j.waitCnd.Signal()
	//					j.waitCnd.L.Unlock()
	//				}
	//
	//				if !j.Active {
	//					//if j was unlocked during wakes handlig (it makes sense - Wake and WakeSync
	//					//use same mutex) - handle possible wakes
	//					if j.cntr > 0 {
	//						//reset state
	//						j.cntr = 0
	//						//handle wakes
	//						...
	//					}
	//					if j.syncStop {
	//						//signal to InternalStop
	//						j.Signal()
	//					}
	//					//Job was stopped
	//					j.Unlock()
	//					return
	//				}
	//			}
	//if isEx == false f should only implement handler for wakes
	isEx bool
	//handler for wake or entire goroutine with loop depends on isEx flag
	f func(*JobCond)
	//Count of wakes received after last handling. Can't be used in handler
	//because resets before handler call. use wakesToHandle
	cntr int
	//Count of wakes received after last handling. Can't be used outside  handler,
	//because value sets before handler call
	WakesToHandle int
	//is some caller waiting for back signal
	isWaiting bool
	//implements back signal
	waitMx *sync.Mutex
	//implements back signal
	waitCnd *sync.Cond
}

//JobChan is wrapper for goroutine that does selection from channel in loop
//goroutine should return if receive  ata from exitChn
type JobChan struct {
	//thread safe Job state
	jobSwitcher
	f       func(*JobChan)
	ExitChn chan struct{}
}

//Start starts job
func (j *JobCond) Start() {
	j.On(j.InternalStart)
}

//Stop stops job
func (j *JobCond) Stop(sync bool) {
	j.Off(func() {
		j.InternalStop(sync)
	})
}

//InternalStart Starts job. Don't call it unless you create new abstraction based on job like Messagequeue
func (j *JobCond) InternalStart() {
	j.isWaiting = false
	j.cntr = 0
	if j.isEx {
		go func() {
			j.Lock()
			j.Signal()
			//j.mx.Unlock()
			j.f(j)

		}()
	} else {
		go func() {
			var isWaiting bool
			j.Lock()
			j.Signal()
			//implement this loop yourself in f function if isEx == true
			for {
				if j.cntr == 0 && j.Active {
					//wait signal
					j.Wait()
				}
				isWaiting = j.isWaiting
				j.isWaiting = false
				j.handleWakes()
				if isWaiting {
					//if WakeSync was called and some caller wait for back signal
					j.signalWakeSync()
				}
				if !j.Active {
					//if j was unlocked during wakes handling (it makes sense - Wake and WakeSync
					//use same mutex) - handle possible wakes
					j.handleWakes()
					if j.syncStop {
						//signal to InternalStop
						j.Signal()
					}
					//Job was stopped
					j.Unlock()
					return
				}
			}
		}()
	}
	j.Wait()
}

//InternalStop Initiates job stop. Don't call it unless you create new abstraction based on job like Messagequeue
func (j *JobCond) InternalStop(sync bool) {
	j.syncStop = sync
	j.Signal()
	if j.syncStop {
		j.Wait()
	}
}

func (j *JobCond) handleWakes() {
	//if Stop and no any Wake happend after last handling j.cntr is 0
	//handle situation when j.cntr > 0 and j.isOn == false in f if you need
	if j.cntr > 0 {
		j.WakesToHandle = j.cntr
		//reset state
		j.cntr = 0
		//call wake handler
		j.f(j)
	}
}

//Wake waiks job goroutine
func (j *JobCond) Wake() {
	j.Lock()
	defer j.Unlock()
	if !j.Active {
		return
	}
	j.InternalWake()
}

//WakeSync waiks job goroutine. Current goroutine stops utill get back signal
func (j *JobCond) WakeSync() {
	j.waitMx.Lock()
	defer j.waitMx.Unlock()

	j.Lock()
	if !j.Active {
		j.Unlock()
		return
	}
	j.isWaiting = true
	j.InternalWake()
	j.Unlock()
	j.waitWakeSync()
}

func (j *JobCond) signalWakeSync() {
	j.waitCnd.L.Lock()
	j.waitCnd.Signal()
	j.waitCnd.L.Unlock()
}

func (j *JobCond) waitWakeSync() {
	j.waitCnd.L.Lock()
	j.waitCnd.Wait()
	j.waitCnd.L.Unlock()
}

//InternalWake does actual waiking of job goroutine. Is not thread safe
func (j *JobCond) InternalWake() {
	j.cntr++
	j.Signal()
}

//Start starts job
func (j *JobChan) Start() {
	j.On(j.InternalStart)
}

//Stop stops job
func (j *JobChan) Stop(sync bool) {
	j.Off(func() {
		j.InternalStop(sync)
	})
}

//InternalStart Starts job. Don't call it unless you create new abstraction based on job like Messagequeue
func (j *JobChan) InternalStart() {
	go func() {
		//signal InternalStart
		j.Lock()
		j.Signal()
		j.Unlock()
		j.f(j)
		if j.syncStop {
			//signal to InternalStop
			j.Lock()
			j.Signal()
			j.Unlock()
		}
	}()
	j.Wait()
}

//InternalStop Initiates job stop. Don't call it unless you create new abstraction based on job like Messagequeue
func (j *JobChan) InternalStop(sync bool) {
	j.syncStop = sync
	j.ExitChn <- struct{}{}
	if j.syncStop {
		j.Wait()
	}
}

//NewJobCond creates new JobCond
func NewJobCond(f func(*JobCond), isEx bool) *JobCond {
	return &JobCond{
		jobSwitcher: newJobSwitcher(),
		isEx:        isEx,
		f:           f,
		waitMx:      new(sync.Mutex),
		waitCnd:     sync.NewCond(new(sync.Mutex))}
}

//NewJobChan creates new JobChan
func NewJobChan(f func(*JobChan)) *JobChan {
	return &JobChan{
		jobSwitcher: newJobSwitcher(),
		f:           f,
		ExitChn:     make(chan struct{})}
}

func newJobSwitcher() jobSwitcher {
	sw := jobSwitcher{Switcher: NewSwitcher()}
	sw.Cond = sync.NewCond(sw.Mutex)
	return sw
}
