package shared

import (
	"sync"
	"testing"
)

var _ MessageSender = new(MessageQueue) //should implement MessageSender

type TestMessage struct {
	text string
}

func TestMessageQueueEx(t *testing.T) {
	l := new(testLogger)
	SetLogger(l)
	wg := new(sync.WaitGroup)

	var mq *MessageQueue
	var msgcntr int
	mq = NewMessageQueueEx(func(m *Message) {
		defer wg.Done()
		msgcntr++
		var data *TestMessage
		if m.Data != nil {
			data = m.Data.(*TestMessage)
		}
		m.Data = 42
		if data == nil {
			l.Log("nil")
			return
		}

		if mq.IsActive() {
			switch m.Code {
			case 42:
				Logf("%d%s", m.Code, data.text)
			case 48:
				Sleep50ms()
				fallthrough
			default:
				Log(data.text)
			}

		} else {
			Log(data.text + "stopping")
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

	mq.SetHandler(0, func(m *Message) {}) // shouldn't raise

	mq.SendSync(0, nil)
	if msgcntr != 0 {
		t.Error("message shouldn't be added before Run")
	}
	mq.Send(0, nil)
	Sleep100ms()
	if msgcntr != 0 {
		t.Error("message shouldn't be added before Run")
	}

	if mq.cntr != 0 {
		t.Error("mq not empty")
	}

	mq.Start(false)

	cntr := 0

	wg.Add(1)
	if mq.SendSync(0, nil) != 42 {
		t.Error("empty message should be added and return result")
	}
	wg.Wait()
	cntr++
	CheckTestLogger(t, l, cntr, "nil")

	wg.Add(1)
	mq.Send(42, &TestMessage{text: "message"})
	wg.Wait()
	cntr++
	CheckTestLogger(t, l, cntr, "42message")

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go mq.Send(1, &TestMessage{text: "message_go"})
	}
	wg.Wait()
	cntr = cntr + 10
	CheckTestLogger(t, l, cntr, "message_go")

	wg.Add(2)
	mq.Send(48, &TestMessage{text: "sync_not_empty_queue_prepare"})
	mq.SendSync(2, &TestMessage{text: "sync_not_empty_queue"})
	//wg.Wait()
	cntr = cntr + 2
	CheckTestLogger(t, l, cntr, "sync_not_empty_queue")

	wg.Add(2)
	mq.Send(3, &TestMessage{text: "exit"})
	Sleep50ms()
	wg.Add(1)
	mq.Send(4, &TestMessage{text: "after_exit"})

	cntr++
	CheckTestLogger(t, l, cntr, "exit")
	wg.Wait()
	Sleep100ms()
	cntr++
	CheckTestLogger(t, l, cntr, "after_exit"+"stopping")

	if mq.SendSync(0, &TestMessage{}) != nil {
		t.Error("message shouldn't be added after stop")
	}
	mq.Stop(false)
}

func TestMessageQueue(t *testing.T) {

	var mq *MessageQueue
	mq = NewMessageQueue()

	var res int

	mq.SendSync(0, nil) // shouldn't raise
	if res != 0 {
		t.Error("wrong res with no handlers")
	}

	mq.SetHandler(0, func(m *Message) {
		res = 0
	})

	mq.SetHandler(1, func(m *Message) {
		res++
	})

	mq.Start(false)

	mq.SendSync(1, nil)
	if res == 0 {
		t.Error("code 1 was not handled")
	}
	mq.SendSync(0, nil)
	if res != 0 {
		t.Error("code 0 was not handled")
	}
	mq.SetHandler(0, nil) //dosen't set if mq.IsOn()
	mq.SendSync(1, nil)
	if res == 0 {
		t.Error("code 1 was not handled")
	}
	mq.SendSync(0, nil)
	if res != 0 {
		t.Error("code 0 should be still handled after mq.SetHandler on live mq")
	}
	mq.Stop(false)
}
func TestMessageQueueStop(t *testing.T) {
	l := new(testLogger)
	SetLogger(l)

	//empty
	wg := new(sync.WaitGroup)
	wg.Add(1)
	mq := NewMessageQueueEx(func(m *Message) {

	})
	mq.Stop(true)

	mq.Start(false)

	if !mq.IsActive() {
		t.Error("Stop empty queue: queue was not started")
	}
	mq.Stop(true)

	if mq.IsActive() {
		t.Error("Stop empty queue: queue was not stopped")
	}

	//not empty
	wg = new(sync.WaitGroup)
	wg.Add(1)
	mq = NewMessageQueueEx(func(m *Message) {
		Sleep10ms()
	})
	mq.Start(false)
	if !mq.IsActive() {
		t.Error("Stop empty queue: queue was not started")
	}
	//handling during aprox. 10* 10ms
	for i := 0; i < 10; i++ {
		mq.Send(0, &TestMessage{text: "message"})
	}
	mq.Stop(false)
	if mq.cntr == 0 {
		t.Error("Stop not empty queue: cntr shoud be great than 0")
	}
	Sleep100ms()
	Sleep10ms()
	if mq.IsActive() {
		t.Error("Stop not empty queue: state should be mqsStoped")
	}
	if mq.cntr != 0 {
		t.Error("Stop not empty queue: cntr shoud be equal 0 after callback")
	}
}
