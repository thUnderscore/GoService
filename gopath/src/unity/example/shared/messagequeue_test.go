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

	mq := NewMessageQueue()

	if mq.Add(&TestMessage{}) {
		t.Error("message shouldn't be added before Run")
	}

	if mq.cnt != 0 {
		t.Error("mq not empty")
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	if !mq.Run(
		func(i interface{}, isRun bool) {
			if i == nil {
				l.Log("nil")
				return
			}
			m := i.(*TestMessage)
			if isRun {
				l.Log(m.text)
			} else {
				l.Log(m.text + "stopping")
			}

			if m.text == "exit" {
				Sleep100ms()
				mq.Stop(func(*MessageQueue) {
					wg.Done()
				})
			}
		}) {
		t.Error("mq isn't ran")
	}

	if mq.Run(nil) {
		t.Error("mq.Run should return false if already ran")
	}

	if !mq.Add(nil) {
		t.Error("empty message be added")
	}
	Sleep100ms()
	CheckTestLogger(t, l, 1, "nil")

	mq.Add(&TestMessage{text: "message"})
	Sleep100ms()
	CheckTestLogger(t, l, 2, "message")

	for i := 0; i < 10; i++ {
		go mq.Add(&TestMessage{text: "message_go"})
	}
	Sleep2s()
	CheckTestLogger(t, l, 12, "message_go")

	mq.Add(&TestMessage{text: "exit"})

	mq.Add(&TestMessage{text: "after_exit"})
	Sleep100ms()
	Sleep10ms()
	CheckTestLogger(t, l, 14, "after_exit"+"stopping")
	wg.Wait()
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
	mq := NewMessageQueue()
	mq.Stop(nil)
	if mq.state != mqsInitial {
		t.Error("mq state should be mqsInitial")
	}
	mq.Run(func(i interface{}, isRun bool) {

	})
	Sleep100ms()
	if mq.state != mqsRun {
		t.Error("Stop empty queue: state should be mqsRun")
	}
	mq.Stop(func(*MessageQueue) {
		wg.Done()
	})
	wg.Wait()
	if mq.state != mqsStoped {
		t.Error("Stop empty queue: state should be mqsStoped")
	}

	//not empty
	wg = new(sync.WaitGroup)
	wg.Add(1)
	mq = NewMessageQueue()
	mq.Run(func(i interface{}, isRun bool) {
		Sleep100ms()
	})
	Sleep100ms()
	if mq.state != mqsRun {
		t.Error("Stop not empty queue: state should be mqsRun")
	}

	for i := 0; i < 10; i++ {
		mq.Add(&TestMessage{text: "message"})
	}
	mq.Stop(func(*MessageQueue) {
		wg.Done()
	})
	if mq.cnt == 0 {
		t.Error("Stop not empty queue: cnt shoud be great than 0")
	}
	wg.Wait()
	if mq.state != mqsStoped {
		t.Error("Stop not empty queue: state should be mqsStoped")
	}
	if mq.cnt != 0 {
		t.Error("Stop not empty queue: cnt shoud be equal 0 after callback")
	}
}
