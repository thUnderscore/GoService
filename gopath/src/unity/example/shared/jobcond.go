package shared

import "sync"

//JobCond is wrapper for goroutine that waits sync.Cond in loop
//goroutine should return if JobCond.isOn == false after Cond.Wait
//Cond could be signaled by Stop or by JobCond.Wake
type JobCond struct {
	//thread safe Job state. Switcher.Mutex used in JobCond.Wake
	*JobSwitcher
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
	//					j.SignalWakeSync()
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
	//                  //j is Locked at the moment
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

//InternalStart Starts job. Don't call it unless you create new abstraction based on job like Messagequeue
func (j *JobCond) InternalStart() {
	j.isWaiting = false
	j.syncStop = false
	j.cntr = 0
	var f func()
	if j.isEx {
		f = func() { j.f(j) }
	} else {
		f = func() {
			var isWaiting bool
			//implement this loop yourself in f function if isEx == true
			for {
				if j.cntr == 0 && j.Active {
					//wait for wake
					j.Wait()
				}
				isWaiting = j.isWaiting
				j.isWaiting = false
				j.handleWakes()
				if isWaiting {
					//if WakeSync was called and some caller wait for back signal
					j.SignalWakeSync()
				}
				if !j.Active {
					//if j was unlocked during wakes handling (it makes sense - Wake and WakeSync
					//use same mutex) - handle possible wakes
					j.handleWakes()
					return
				}
			}
		}
	}
	j.RunGoroutine(f, true)
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

//SignalWakeSync signals caller tah has called WakeSync
func (j *JobCond) SignalWakeSync() {
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

//NewJobCond creates new JobCond
func NewJobCond(f func(*JobCond), isEx bool) *JobCond {
	j := &JobCond{
		JobSwitcher: NewJobSwitcher(),
		isEx:        isEx,
		f:           f,
		waitMx:      new(sync.Mutex),
		waitCnd:     sync.NewCond(new(sync.Mutex)),
	}
	j.OnStart = j.InternalStart
	j.OnStop = j.Signal
	return j
}
