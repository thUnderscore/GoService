package shared

import "testing"
import "sync"

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
			case <-j.ExitChn:
				return
			}

		}
	})

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
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			Log("before wake")
			cnd.L.Lock()
			sgnl <- 1
			cnd.Wait()
			cnd.L.Unlock()
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

	j.Stop(true)
	if j.IsActive() {
		t.Error("Job should be stopped")
	}

	j.Start(true)
	if j.IsActive() {
		t.Error("Restart: Job should not be started (once == true)")
	}
	j.Start(false)
	//wg.Add(1)

	cnd.L.Lock()
	sgnl <- 2
	cnd.Wait()
	cnd.L.Unlock()

	//wg.Wait()
	cnt = cnt + 1
	CheckTestLogger(t, l, cnt, "handle JobChan")
	if cntr != 12 {
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
