package shared

import "time"

func Sleep100ms() {
	time.Sleep(100 * time.Millisecond)
}

func Sleep50ms() {
	time.Sleep(50 * time.Millisecond)
}

func Sleep10ms() {
	time.Sleep(10 * time.Millisecond)
}

func Sleep1ms() {
	time.Sleep(1 * time.Millisecond)
}

func Sleep2s() {
	time.Sleep(2000 * time.Millisecond)
}

func Sleep1s() {
	time.Sleep(1000 * time.Millisecond)
}
