package shared

import "sync"
import "fmt"

var logger Logger
var mu = new(sync.RWMutex)

//Logger interface that wraps basic Log method
type Logger interface {
	//write string to log
	Log(str string)
}

//Log sends string to log
func Log(str string) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.Log(str)
	} else {
		fmt.Println(str)
	}
}

//Logf formats and sends string to log
func Logf(format string, a ...interface{}) {
	Log(fmt.Sprintf(format, a...))
}

//SetLogger sets global object that provides Logger interface
func SetLogger(lgr Logger) {
	mu.Lock()
	defer mu.Unlock()
	logger = lgr
}

//ResetLogger sets global logger to nil
func ResetLogger() {
	SetLogger(nil)
}

//HasLogger returns if global logger is set
func HasLogger() bool {
	mu.Lock()
	defer mu.Unlock()
	return logger != nil
}
