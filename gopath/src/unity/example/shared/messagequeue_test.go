package shared

import (
	"sync"
	"testing"
)

type TestMessage struct {
	text string
}

func TestMessageQueue(t *testing.T) {
	l := new(testLogger)
	SetLogger(l)
	wg := new(sync.WaitGroup)

	var mq *MessageQueue
	mq = NewMessageQueue(func(i interface{}, isOn bool) {
		if i == nil {
			l.Log("nil")
			return
		}
		m := i.(*TestMessage)
		if isOn {
			l.Log(m.text)
		} else {
			l.Log(m.text + "stopping")
		}

		if m.text == "exit" {
			Sleep100ms()
			go func() {
				mq.Stop(true)
				wg.Done()
			}()
			Sleep50ms()
		}
	})

	if mq.Add(&TestMessage{}) {
		t.Error("message shouldn't be added before Run")
	}

	if mq.cntr != 0 {
		t.Error("mq not empty")
	}

	wg.Add(1)
	mq.Start()

	cntr := 0

	if !mq.Add(nil) {
		t.Error("empty message should be added")
	}
	Sleep100ms()
	cntr++
	CheckTestLogger(t, l, cntr, "nil")

	mq.Add(&TestMessage{text: "message"})
	Sleep100ms()
	cntr++
	CheckTestLogger(t, l, cntr, "message")

	for i := 0; i < 10; i++ {
		go mq.Add(&TestMessage{text: "message_go"})
	}
	Sleep100ms()
	cntr = cntr + 10
	CheckTestLogger(t, l, cntr, "message_go")

	mq.Add(&TestMessage{text: "exit"})
	Sleep50ms()
	mq.Add(&TestMessage{text: "after_exit"})
	cntr++
	CheckTestLogger(t, l, cntr, "exit")
	wg.Wait()
	Sleep100ms()
	cntr++
	CheckTestLogger(t, l, cntr, "after_exit"+"stopping")

	if mq.Add(&TestMessage{}) {
		t.Error("message shouldn't be added after stop")
	}
}

func TestMessageQueueStop(t *testing.T) {
	l := new(testLogger)
	SetLogger(l)

	//empty
	wg := new(sync.WaitGroup)
	wg.Add(1)
	mq := NewMessageQueue(func(i interface{}, isOn bool) {

	})
	mq.Stop(true)

	mq.Start()

	if !mq.isOn() {
		t.Error("Stop empty queue: queue was not started")
	}
	mq.Stop(true)

	if mq.isOn() {
		t.Error("Stop empty queue: queue was not stopped")
	}

	//not empty
	wg = new(sync.WaitGroup)
	wg.Add(1)
	mq = NewMessageQueue(func(i interface{}, isRun bool) {
		Sleep50ms()
	})
	mq.Start()
	if !mq.isOn() {
		t.Error("Stop empty queue: queue was not started")
	}

	for i := 0; i < 10; i++ {
		mq.Add(&TestMessage{text: "message"})
	}
	Sleep100ms()
	for i := 0; i < 10; i++ {
		mq.Add(&TestMessage{text: "message"})
	}
	mq.Stop(false)
	if mq.cntr == 0 {
		t.Error("Stop not empty queue: cntr shoud be great than 0")
	}
	Sleep1s()
	if mq.isOn() {
		t.Error("Stop not empty queue: state should be mqsStoped")
	}
	if mq.cntr != 0 {
		t.Error("Stop not empty queue: cntr shoud be equal 0 after callback")
	}
}
