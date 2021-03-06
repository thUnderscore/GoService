package shared

import "testing"
import "sync"

func TestJobCondSmpl(t *testing.T) {
	ResetLogger()
	l := new(testLogger)
	SetLogger(l)

	wg := new(sync.WaitGroup)

	cntr := 0
	j := NewJobCond(func(j *JobCond) {
		cntr = cntr + j.WakesToHandle
		for i := 0; i < j.WakesToHandle; i++ {
			Log("handle JobCond")
			wg.Done()
		}
	}, false)
	var _ Job = j //should implement Job
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			j.Start(true)
			wg.Done()
		}()
	}
	j.Start(true)
	wg.Wait()
	if !j.IsActive() {
		t.Error("Job should be started")
	}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			Log("before wake")
			j.WakeSync()
			Log("after wake")
			wg.Done()
		}()
	}
	wg.Wait()
	cnt := 30
	CheckTestLogger(t, l, cnt, "after wake")
	if cntr != 10 {
		t.Error("cntr should be ", 10, "not", cntr)
	}

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			j.Wake()
		}()
	}
	wg.Wait()
	if cntr != 20 {
		t.Error("wake : cntr should be ", 20, "not", cntr)
	}
	/*
		for i := 0; i < 10; i++ {
			go j.Stop(false)
		}
	*/
	j.Stop(true)
	j.Wake()
	if j.cntr != 0 {
		t.Error("wake on stopped job should be ignorred ")
	}
	j.WakeSync() //shouldn't hang
	if j.IsActive() {
		t.Error("Job should be stopped")
	}

	j.Start(true)
	if j.IsActive() {
		t.Error("Restart: Job should not be started (once == true)")
	}
	j.Start(false)
	if !j.IsActive() {
		t.Error("Restart: Job should be started")
	}
	wg.Add(1)
	j.Wake()
	wg.Wait()
	cnt = cnt + 1
	if cntr != 21 {
		t.Error("Wake was not processed")
	}

	for i := 0; i < 10; i++ {
		go j.Stop(false)
	}
	Sleep100ms()
	if j.IsActive() {
		t.Error("Restart: Job should be stopped")
	}

}

func TestJobCondEx(t *testing.T) {
	ResetLogger()
	l := new(testLogger)
	SetLogger(l)

	wg := new(sync.WaitGroup)

	cntr := 0
	j := NewJobCond(func(j *JobCond) {

		for {
			if j.cntr == 0 && j.Active {
				//wait signal
				j.Wait()
			}
			//if Stop and no any Wake happend after last handling j.cntr is 0
			//handle situation when j.cntr > 0 and j.isOn == false in f if you need
			if j.cntr > 0 {
				//handle gignals
				j.WakesToHandle = j.cntr
				cntr = cntr + j.WakesToHandle
				for i := 0; i < j.WakesToHandle; i++ {
					Log("handle JobCondEx")
					wg.Done()
				}

				//reset wakes counter
				j.cntr = 0
			}

			if j.isWaiting {
				//if WakeSync was called and some caller wait for back signal
				j.waitCnd.L.Lock()
				j.waitCnd.Signal()
				j.waitCnd.L.Unlock()
			}

			if !j.Active {
				//if j.mx was unlocked during wakes handlig (it makes sense - Wake and WakeSync
				//use same mutex) - handle possible wakes
				if j.cntr > 0 {
					//handle gignals
					j.WakesToHandle = j.cntr
					cntr = cntr + j.WakesToHandle
					Log("handle JobCondEx")
					wg.Done()
					//reset wakes counter
					j.cntr = 0
				}
				return
			}
		}

	}, true)
	var _ Job = j //should implement Job

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			j.Start(true)
			wg.Done()
		}()
	}
	j.Start(true)
	wg.Wait()
	if !j.IsActive() {
		t.Error("Job should be started")
	}

	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			Log("before wake")
			j.WakeSync()
			Log("after wake")
			wg.Done()
		}()
	}
	wg.Wait()
	cnt := 30
	CheckTestLogger(t, l, cnt, "after wake")

	if cntr != 10 {
		t.Error("cntr should be ", 10, "not", cntr)
	}

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			j.Wake()
		}()
	}
	wg.Wait()
	if cntr != 20 {
		t.Error("wake: cntr should be ", 20, "not", cntr)
	}

	j.Stop(true)
	/*
		for i := 0; i < 10; i++ {
			go j.Stop(false)
		}
	*/
	if j.IsActive() {
		t.Error("Job should be stopped")
	}

	j.Start(true)
	if j.IsActive() {
		t.Error("Restart: Job should not be started (once == true)")
	}
	j.Start(false)
	if !j.IsActive() {
		t.Error("Restart: Job should be started")
	}
	wg.Add(1)
	j.Wake()
	wg.Wait()
	cnt = cnt + 1
	if cntr != 21 {
		t.Error("Wake was not processed")
	}

	for i := 0; i < 10; i++ {
		go j.Stop(false)
	}
	Sleep100ms()
	if j.IsActive() {
		t.Error("Restart: Job should be stopped")
	}
}
