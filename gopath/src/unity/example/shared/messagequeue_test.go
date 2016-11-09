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
	mq = NewMessageQueue(func(m *Message, isOn bool) {
		var data *TestMessage
		if m.data != nil {
			data = m.data.(*TestMessage)
		}
		m.data = 42
		if data == nil {
			l.Log("nil")
			return
		}

		if isOn {
			l.Log(data.text)
		} else {
			l.Log(data.text + "stopping")
		}

		if data.text == "exit" {
			Sleep100ms()
			go func() {
				mq.Stop(true)
				wg.Done()
			}()
			Sleep50ms()
		}
	})

	if mq.Send(0, nil, true) != nil {
		t.Error("message shouldn't be added before Run")
	}

	if mq.cntr != 0 {
		t.Error("mq not empty")
	}

	wg.Add(1)
	mq.Start()

	cntr := 0

	if mq.Send(0, nil, true) != 42 {
		t.Error("empty message should be added and return result")
	}
	Sleep100ms()
	cntr++
	CheckTestLogger(t, l, cntr, "nil")

	mq.Send(0, &TestMessage{text: "message"}, false)
	Sleep100ms()
	cntr++
	CheckTestLogger(t, l, cntr, "message")

	for i := 0; i < 10; i++ {
		go mq.Send(0, &TestMessage{text: "message_go"}, false)
	}
	Sleep100ms()
	cntr = cntr + 10
	CheckTestLogger(t, l, cntr, "message_go")

	mq.Send(0, &TestMessage{text: "exit"}, false)
	Sleep50ms()
	mq.Send(0, &TestMessage{text: "after_exit"}, false)
	cntr++
	CheckTestLogger(t, l, cntr, "exit")
	wg.Wait()
	Sleep100ms()
	cntr++
	CheckTestLogger(t, l, cntr, "after_exit"+"stopping")

	if mq.Send(0, &TestMessage{}, true) != nil {
		t.Error("message shouldn't be added after stop")
	}
}

func TestMessageQueueStop(t *testing.T) {
	l := new(testLogger)
	SetLogger(l)

	//empty
	wg := new(sync.WaitGroup)
	wg.Add(1)
	mq := NewMessageQueue(func(m *Message, isOn bool) {

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
	mq = NewMessageQueue(func(m *Message, isRun bool) {
		Sleep50ms()
	})
	mq.Start()
	if !mq.isOn() {
		t.Error("Stop empty queue: queue was not started")
	}

	for i := 0; i < 10; i++ {
		mq.Send(0, &TestMessage{text: "message"}, false)
	}
	Sleep100ms()
	for i := 0; i < 10; i++ {
		mq.Send(0, &TestMessage{text: "message"}, false)
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
