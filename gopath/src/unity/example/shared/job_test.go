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
		cntr = cntr + j.wakesToHandle
		Log("handle JobCond")
		wg.Done()
	}, false)

	for i := 0; i < 10; i++ {
		go j.Start()
	}
	j.Start()
	if !j.isOn() {
		t.Error("Job should be started")
	}
	wg.Add(11)
	for i := 0; i < 10; i++ {
		go func() {
			Log("before wake")
			j.WakeSync()
			Log("after wake")
		}()
	}
	Sleep100ms()
	cnt := 30
	CheckTestLogger(t, l, cnt, "after wake")
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
			j.Wake()
		}()
	}
	Sleep1s()
	//wg.Wait()
	if cntr != 20 {
		t.Error("wake : cntr should be ", 20, "not", cntr)
	}

	for i := 0; i < 10; i++ {
		go j.Stop(false)
	}
	j.Stop(true)
	if j.isOn() {
		t.Error("Job should be stopped")
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
			//wait signal
			j.cnd.Wait()
			//if Stop and no any Wake happend after last handling j.cntr is 0
			//handle situation when j.cntr > 0 and j.isOn == false in f if you need
			if j.cntr > 0 {
				//handle gignals
				j.wakesToHandle = j.cntr
				cntr = cntr + j.wakesToHandle
				Log("handle JobCondEx")
				wg.Done()

				//reset wakes counter
				j.cntr = 0
			}

			if j.isWaiting {
				//if WakeSync was called and some caller wait for back signal
				j.waitCnd.L.Lock()
				j.waitCnd.Signal()
				j.waitCnd.L.Unlock()
			}

			if !j.on {
				//if j.mx was unlocked during wakes handlig (it makes sense - Wake and WakeSync
				//use same mutex) - handle possible wakes
				if j.cntr > 0 {
					//handle gignals
					j.wakesToHandle = j.cntr
					cntr = cntr + j.wakesToHandle
					Log("handle JobCondEx")
					wg.Done()
					//reset wakes counter
					j.cntr = 0
				}
				if j.syncStop {
					//signal to InternalStop
					j.cnd.Signal()
				}
				//Job was stopped
				j.mx.Unlock()
				return
			}
		}

	}, true)
	for i := 0; i < 10; i++ {
		go j.Start()
	}
	j.Start()
	if !j.isOn() {
		t.Error("Job should be started")
	}

	wg.Add(11)
	for i := 0; i < 10; i++ {
		go func() {
			Log("before wake")
			j.WakeSync()
			Log("after wake")
		}()
	}
	Sleep100ms()
	cnt := 30
	CheckTestLogger(t, l, cnt, "after wake")
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
			j.Wake()
		}()
	}
	Sleep1s()
	//wg.Wait()
	if cntr != 20 {
		t.Error("wake: cntr should be ", 20, "not", cntr)
	}

	for i := 0; i < 10; i++ {
		go j.Stop(false)
	}
	j.Stop(true)
	if j.isOn() {
		t.Error("Job should be stopped")
	}
}

func TestJobChan(t *testing.T) {
	ResetLogger()
	l := new(testLogger)
	SetLogger(l)

	wg := new(sync.WaitGroup)

	cntr := 0

	cnd := sync.NewCond(new(sync.Mutex))

	sgnl := make(chan int, 10)
	j := NewJobChan(func(j *JobChan) {
		for {
			select {
			case b := <-sgnl:
				cntr = cntr + b
				Log("handle JobChan")
				cnd.L.Lock()
				cnd.Signal()
				cnd.L.Unlock()
			case <-j.exitChn:
				return
			}

		}
	})

	for i := 0; i < 10; i++ {
		go j.Start()
	}
	j.Start()
	if !j.isOn() {
		t.Error("Job should be started")
	}
	wg.Add(11)
	for i := 0; i < 10; i++ {
		go func() {
			Log("before wake")
			sgnl <- 1
			cnd.L.Lock()
			cnd.Wait()
			cnd.L.Unlock()
			Log("after wake")
			wg.Done()
		}()
	}
	Sleep100ms()
	cnt := 30
	CheckTestLogger(t, l, cnt, "after wake")
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

	for i := 0; i < 10; i++ {
		go j.Stop(false)
	}
	j.Stop(true)
	if j.isOn() {
		t.Error("Job should be stopped")
	}

}
