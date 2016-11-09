package shared

import "sync"

type queueItem struct {
	next *queueItem
	data interface{}
}

//MessageQueue structure represent message queue
type MessageQueue struct {
	JobCond
	//head of messages list
	head *queueItem
	//tail of messages list
	tail    *queueItem
	handler func(interface{}, bool)
}

//queueItem pool
var itemFree = sync.Pool{
	New: func() interface{} { return new(queueItem) },
}

// newPrinter allocates a new pp struct or grabs a cached one.
func newQueueItem(data interface{}) *queueItem {
	i := itemFree.Get().(*queueItem)
	i.data = data
	return i
}

// free saves used pp structs in ppFree; avoids an allocation per invocation.
func (i *queueItem) free() {
	i.next = nil
	i.data = nil
	itemFree.Put(i)
}

//Add adds message to queue. If queue is not running
//msg is not added and func returns true
func (q *MessageQueue) Add(msg interface{}) bool {
	q.mx.Lock()
	defer q.mx.Unlock()
	if !q.on {
		return false
	}

	i := newQueueItem(msg)
	//msg.setNext(nil)
	if q.head == nil {
		q.head = i
	} else {
		q.tail.next = i
	}
	q.tail = i
	q.InternalWake()
	return true
}

//NewMessageQueue MessageQueue constructor
func NewMessageQueue(handler func(interface{}, bool)) *MessageQueue {
	mq := &MessageQueue{handler: handler}
	mq.JobCond = *NewJobCond(func(j *JobCond) {
		curr := mq.head
		isOn := mq.on
		mq.head = nil
		mq.tail = nil
		var t *queueItem
		j.mx.Unlock()
		defer j.mx.Lock()
		for curr != nil {
			handler(curr.data, isOn)
			t = curr
			curr = curr.next
			t.free()

		}
	}, false)
	return mq
}
