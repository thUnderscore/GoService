package shared

import (
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

type testLogger struct {
	m    sync.Mutex
	cnt  int
	last string
}

func (lgr *testLogger) Log(str string) {
	lgr.m.Lock()
	defer lgr.m.Unlock()
	lgr.last = str
	lgr.cnt++
}

func TestLog(t *testing.T) {
	Log("0")
	for i := 0; i < 10; i++ {
		l := new(testLogger)
		SetLogger(l)
		Log("a")
		CheckTestLogger(t, l, 1, "a")
		Log("b")
		CheckTestLogger(t, l, 2, "b")
		Logf("c %v", 3)
		CheckTestLogger(t, l, 3, "c 3")
		wg := new(sync.WaitGroup)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Log("10")
			}()
		}

		wg.Wait()
		if !HasLogger() {
			t.Errorf("has no logger")
		}
		CheckTestLogger(t, l, 13, "10")
		ResetLogger()
		if HasLogger() {
			t.Errorf("has logger")
		}
	}

}

func CheckTestLogger(t *testing.T, l *testLogger, cnt int, last string) {
	if l.cnt != cnt || l.last != last {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			t.Errorf("Expected %v:%v actual %v:%v at file: %v line: %v", cnt, last, l.cnt, l.last, filepath.Base(file), line)
		} else {
			t.Errorf("Expected %v:%v actual %v:%v", cnt, last, l.cnt, l.last)
		}
	}
}
