package shared

//JobChan is wrapper for goroutine that does selection from channel in loop
//goroutine should return if receive  ata from exitChn
type JobChan struct {
	//thread safe Job state
	*JobSwitcher
	f       func(*JobChan)
	ExitChn chan struct{}
}

//InternalStart Starts job. Don't call it unless you create new abstraction based on job like Messagequeue
func (j *JobChan) InternalStart() {
	j.RunGoroutine(func() { j.f(j) }, false)

}

//InternalStop Initiates job stop. Don't call it unless you create new abstraction based on job like Messagequeue
func (j *JobChan) InternalStop() {
	j.ExitChn <- struct{}{}
}

//NewJobChan creates new JobChan
func NewJobChan(f func(*JobChan)) *JobChan {
	j := &JobChan{
		JobSwitcher: NewJobSwitcher(),
		f:           f,
		ExitChn:     make(chan struct{}),
	}
	j.OnStart = j.InternalStart
	j.OnStop = j.InternalStop
	return j
}
