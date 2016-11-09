package shared

//MessageQueue structure represent message queue
type MessageQueue struct {
	JobCond
	//head of messages list
	head *Message
	//tail of messages list
	tail    *Message
	handler func(data *Message, isOn bool)
}

//Send adds message to queue. If queue is not running
//message is not added and func returns true
func (q *MessageQueue) Send(code MessageCode, data interface{}, sync bool) interface{} {

	q.mx.Lock()
	if !q.on {
		q.mx.Unlock()
		return nil
	}

	m := newMessage(code, data, sync)
	//data.setNext(nil)
	if q.head == nil {
		q.head = m
	} else {
		q.tail.next = m
	}
	q.tail = m
	q.InternalWake()
	q.mx.Unlock()
	if !sync {
		return nil
	}
	return m.wait()
}

//NewMessageQueue MessageQueue constructor
func NewMessageQueue(handler func(*Message, bool)) *MessageQueue {
	mq := &MessageQueue{handler: handler}
	mq.JobCond = *NewJobCond(func(j *JobCond) {
		curr := mq.head
		isOn := mq.on
		mq.head = nil
		mq.tail = nil
		var m *Message
		j.mx.Unlock()
		for curr != nil {
			handler(curr, isOn)
			m = curr
			curr = curr.next
			m.handled()
		}
		j.mx.Lock()
	}, false)
	return mq
}
