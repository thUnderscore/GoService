package shared

//MessageQueue structure represent message queue
type MessageQueue struct {
	*JobCond
	*MessageHandler
	//head of messages list
	head *Message
	//tail of messages list
	tail *Message
}

//Send adds message to queue. If queue is not running
//message is not added and func returns true
func (q *MessageQueue) Send(code MessageCode, data interface{}) {
	q.Lock()
	if !q.Active {
		q.Unlock()
		return
	}
	m := NewMessage(code, data, false)
	if q.head == nil {
		q.head = m
	} else {
		q.tail.Next = m
	}
	q.tail = m
	q.InternalWake()
	q.Unlock()
}

//SendSync adds message to queue and wait result. If queue is not running
//message is not added and func returns true
func (q *MessageQueue) SendSync(code MessageCode, data interface{}) interface{} {
	q.Lock()
	if !q.Active {
		q.Unlock()
		return nil
	}
	m := NewMessage(code, data, true)
	if q.head == nil {
		q.head = m
	} else {
		q.tail.Next = m
	}
	q.tail = m
	q.InternalWake()
	q.Unlock()
	return m.Wait()
}

//NewMessageQueueEx MessageQueue constructor. handler calls for every message in queue
//SetHandler does not work if you create MessageQueue using this function
func NewMessageQueueEx(handler func(*Message)) *MessageQueue {
	return newMessageQueue(handler)
}

//NewMessageQueue MessageQueue constructor. MessageQueue uses embeded MessageHandler
//to handle messages. Use SetHandler to assign handler to  code
func NewMessageQueue() *MessageQueue {
	handler := NewMessageHandler()
	mq := newMessageQueue(handler.Handle)
	mq.MessageHandler = handler
	return mq
}

func newMessageQueue(handler func(*Message)) *MessageQueue {
	mq := &MessageQueue{}
	mq.JobCond = NewJobCond(func(j *JobCond) {
		m := mq.head
		mq.head = nil
		mq.tail = nil
		j.Unlock()
		m.Handle(handler)
		j.Lock()
	}, false)

	return mq
}

//SetHandler sets handler associated with message code. If f is nil handler removes
//DOesn't work if you create MessageQueue using NewMessageQueueEx
func (q *MessageQueue) SetHandler(code MessageCode, f func(m *Message)) {
	q.Lock()
	defer q.Unlock()
	if q.MessageHandler == nil || q.Active {
		return
	}
	q.MessageHandler.SetHandler(code, f)
}
