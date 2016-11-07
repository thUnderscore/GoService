package shared

import "testing"
import "sync"

func TestJobCond(t *testing.T) {
	ResetLogger()
	l := new(testLogger)
	SetLogger(l)

	wg := new(sync.WaitGroup)

	cntr := 0
	j := NewJobCond(func(j *JobCond) {
		cntr = cntr + j.signals
		Log("handle")
		wg.Done()
	}, false)

	for i := 0; i < 10; i++ {
		go j.Start()
	}
	if !j.isOn() {
		t.Error("Job should be started")
	}
	wg.Add(11)
	for i := 0; i < 10; i++ {
		go func() {
			Log("before signal")
			j.Signal(true)
			Log("after signal")
		}()
	}
	Sleep100ms()
	cnt := 30
	CheckTestLogger(t, l, cnt, "after signal")
	go func() {
		Sleep1s()
		Log("defered wait done")
		wg.Done()
	}()
	wg.Wait()
	if cntr != 10 {
		t.Error("cntr should be ", 10, "not", cntr)
	}
	cnt = cnt + 1
	CheckTestLogger(t, l, cnt, "defered wait done")

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			j.Signal(false)
		}()
	}
	Sleep1s()
	//wg.Wait()
	if cntr != 20 {
		t.Error("Signal with no wait: cntr should be ", 20, "not", cntr)
	}

	for i := 0; i < 10; i++ {
		go j.Stop()
	}

	Sleep100ms()
	if j.isOn() {
		t.Error("Job should be stopped")
	}

}

func TestSvc(t *testing.T) {

}
