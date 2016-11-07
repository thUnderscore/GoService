package shared

import "sync"

//Job is interface that describe wrapper for goroutine
type Job interface {
	Start()
	Stop()
	isOn()
}

//JobCond is wrapper for goroutine that waits sync.Cond in loop
//goroutine should return if JobCond.isOn == false after Cond.Wait
//Cond could be signaled by Stop or by JobCond.Signal
type JobCond struct {
	//thread safe Job state. Switcher.Mutex used in JobCond.Signal
	Switcher
	//if isEx == true JobCond.f should implement for loop like in JobCond.startFunc
	//implement this loop yourself in f function if isEx == true
	//			for {
	//				//wait signal
	//				cnd.Wait()
	//				//if Stop and no any Signal happend after last handling j.signals is 0
	//				//handle situation when j.signals > 0 and j.isOn == false in f if you need
	//				if j.signals > 0 {
	//					//handle gignals
	//					...
	//					//reset signals counter
	//					j.signals = 0
	//				}
	//
	//				if j.isWaiting {
	//					//if Signal(true) was called and some caller wait for back signal
	//					wCnd.L.Lock()
	//					wCnd.Signal()
	//					wCnd.L.Unlock()
	//				}
	//
	//				if !j.on {
	//					//Job was stopped
	//					l.Unlock()
	//					return
	//				}
	//			}
	//if isEx == false f should only implement handler for signals
	isEx bool
	//handler for signal or entire goroutine with loop depends on isEx flag
	f func(*JobCond)
	//Cond used for goroutine loop  implementation
	cnd *sync.Cond
	//Count of signals received after last handling
	signals int
	//is some caller waiting for back signal
	isWaiting bool
	//implements back signal
	waitMx *sync.Mutex
	//implements back signal
	waitCnd *sync.Cond
}

//JobChan is wrapper for goroutine that does selection from channel in loop
//goroutine should return if receive  ata from exChn
type JobChan struct {
	Switcher
	f     func(*JobChan)
	exChn chan struct{}
}

//Start starts job
func (j *JobCond) Start() {
	j.On(j.startFunc)
}

//Stop stops job
func (j *JobCond) Stop() {
	j.Off(j.stopFunc)
}

func (j *JobCond) startFunc() {
	if j.isEx {
		go func() {
			j.cnd.L.Lock()
			j.cnd.Signal()
			j.cnd.L.Unlock()
			j.f(j)
		}()
	} else {
		go func() {

			f := j.f
			l := j.mx
			cnd := j.cnd
			wCnd := j.waitCnd
			l.Lock()
			cnd.Signal()
			//implement this loop yourself in f function if isEx == true
			for {
				//wait signal
				cnd.Wait()
				//if Stop and no any Signal happend after last handling j.signals is 0
				//handle situation when j.signals > 0 and j.isOn == false in f if you need
				if j.signals > 0 {
					//call signal handler
					f(j)
					//reset signals counter
					j.signals = 0
				}

				if j.isWaiting {
					//if Signal(true) was called and some caller wait for back signal
					wCnd.L.Lock()
					wCnd.Signal()
					wCnd.L.Unlock()
				}

				if !j.on {
					//Job was stopped
					l.Unlock()
					return
				}
			}
		}()
	}
	j.cnd.Wait()
}

func (j *JobCond) stopFunc() {
	j.cnd.Signal()
}

//Signal waiks job goroutine. if wait == true current goroutine stops utill get back signal
func (j *JobCond) Signal(wait bool) {

	if wait {
		j.waitMx.Lock()
		defer j.waitMx.Unlock()
	}

	j.mx.Lock()
	j.isWaiting = j.isWaiting || wait
	j.signals++
	j.cnd.Signal()
	j.mx.Unlock()

	if wait {
		j.waitCnd.L.Lock()
		j.waitCnd.Wait()
		j.waitCnd.L.Unlock()
	}
}

//Start starts job
func (j *JobChan) Start() {
	j.On(j.startFunc)
}

//Stop stops job
func (j *JobChan) Stop() {
	j.Off(j.stopFunc)
}

func (j *JobChan) startFunc() {
	go j.f(j)
}
func (j *JobChan) stopFunc() {
	j.exChn <- struct{}{}
}

//NewJobCond creates new JobCond
func NewJobCond(f func(*JobCond), isEx bool) *JobCond {

	sw := NewSwitcher()
	return &JobCond{
		Switcher: sw,
		isEx:     isEx,
		f:        f,
		cnd:      sync.NewCond(sw.mx),
		waitMx:   new(sync.Mutex),
		waitCnd:  sync.NewCond(new(sync.Mutex))}
}

//NewJobChan creates new JobChan
func NewJobChan(f func(*JobChan)) *JobChan {
	return &JobChan{
		Switcher: NewSwitcher(),
		f:        f,
		exChn:    make(chan struct{})}
}
