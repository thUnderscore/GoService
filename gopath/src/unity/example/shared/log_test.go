package shared

import (
	"bufio"
	"os"
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
	ResetLogger()

	oldStdout := os.Stdout
	readFile, writeFile, err := os.Pipe()
	if err != nil {
		t.Error("cant create pipe", err)
	}
	os.Stdout = writeFile
	var s string
	go func() {
		scanner := bufio.NewScanner(readFile)
		for scanner.Scan() {
			line := scanner.Text()
			s = s + line
		}
	}()
	Sleep100ms()
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
		Log("11")
		CheckTestLogger(t, l, 13, "10")
		if HasLogger() {
			t.Errorf("has logger")
		}
	}

	writeFile.Close()
	os.Stdout = oldStdout
	if s != "11111111111111111111" {
		t.Error("Reseted log: ", s)
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
