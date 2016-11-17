package shared

//JobCustom is wrapper for goroutine that does selection from channel in loop
//goroutine should return if receive  ata from exitChn
type JobCustom struct {
	//thread safe Job state
	*JobSwitcher
	f func()
}

//InternalStart Starts job. Don't call it unless you create new abstraction based on job like Messagequeue
func (j *JobCustom) InternalStart() {
	j.RunGoroutine(j.f, false)

}

//NewJobCustom creates new JobCustom
func NewJobCustom(f func()) *JobCustom {
	j := &JobCustom{
		JobSwitcher: NewJobSwitcher(),
		f:           f,
	}
	j.OnStart = j.InternalStart
	return j
}
