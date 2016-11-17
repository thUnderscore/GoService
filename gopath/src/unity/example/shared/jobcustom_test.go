package shared

import "testing"
import "sync"

func TestJobCustom(t *testing.T) {
	ResetLogger()
	l := new(testLogger)
	SetLogger(l)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	cntr := 0

	var j *JobCustom
	j = NewJobCustom(func() {
		for j.IsActive() {
			cntr++
			if cntr == 1 {
				wg.Done()
			}
			Sleep100ms()
		}
	})

	var _ Job = j //should implement Job
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			j.Start(true)
			wg.Done()
		}()
	}
	j.Start(true)
	if !j.IsActive() {
		t.Error("Job should be started")
	}
	wg.Wait()
	if cntr != 1 {
		t.Error("cntr should be ", 1, "not", cntr)
	}

	j.Stop(true)
	if j.IsActive() {
		t.Error("Job should be stopped")
	}

	j.Start(true)
	if j.IsActive() {
		t.Error("Job should be restarted")
	}
	cntr = 0
	wg.Add(1)
	j.Start(false)
	if !j.IsActive() {
		t.Error("Restart: Job should be started")
	}
	wg.Wait()
	if cntr != 1 {
		t.Error("cntr should be ", 1, "not", cntr)
	}

	for i := 0; i < 10; i++ {
		go j.Stop(false)
	}
	Sleep100ms()
	Sleep50ms()
	if j.IsActive() {
		t.Error("Restart: Job should be stopped")
	}

}
