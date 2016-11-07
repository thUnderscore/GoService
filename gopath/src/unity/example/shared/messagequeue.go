package shared

import "sync"

//QueueState state of MessageQueue
type QueueState byte

const (
	//Initial state of MessageQueue
	mqsInitial QueueState = iota
	//State of MessageQueue after Run before Stop
	mqsRun
	//State of MessageQueue after Stop
	mqsStoped
)

type queueItem struct {
	next *queueItem
	data interface{}
}

//MessageQueue structure represent message queue
type MessageQueue struct {
	//head of messages list
	head *queueItem
	//tail of messages list
	tail *queueItem
	//Count of messages in queue
	cnt int
	//mutex that protects structure fields
	mx *sync.Mutex
	//Cond that singnals about stop\add events
	cnd *sync.Cond
	//state of meggase queue. See constants
	state QueueState
	//func to be called after queue stop
	stopClbck func(*MessageQueue)
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

//NewMessageQueue MessageQueue constructor
func NewMessageQueue() *MessageQueue {
	mx := &sync.Mutex{}
	return &MessageQueue{mx: mx, cnd: sync.NewCond(mx), state: mqsInitial}
}

//Add adds message to queue. If queue is not running
//msg is not added and func returns true
func (q *MessageQueue) Add(msg interface{}) bool {
	/*if msg == nil {
		return false
	}*/
	q.mx.Lock()
	defer q.mx.Unlock()
	if q.state != mqsRun {
		return false
	}
	q.cnt++
	i := newQueueItem(msg)
	//msg.setNext(nil)
	if q.head == nil {
		q.head = i
		q.tail = i
		q.cnd.Signal()
	} else {
		q.tail.next = i
		q.tail = i
	}

	return true
}

//Stop stops message queue. If there are messages in queue they whould be handled
//with second param == true
func (q *MessageQueue) Stop(stopClbck func(*MessageQueue)) {
	q.mx.Lock()
	defer q.mx.Unlock()
	if q.state != mqsRun {
		return
	}
	q.state = mqsStoped
	q.stopClbck = stopClbck
	if q.head == nil {
		q.cnd.Signal()
	}
}

//Run starts message queue. Handler is func that handles each message
//first param is message, second - is queue running. if queue is running
//or stopped Run returns false
func (q *MessageQueue) Run(handler func(interface{}, bool)) bool {
	q.mx.Lock()
	if q.state != mqsInitial {
		q.mx.Unlock()
		return false
	}
	q.state = mqsRun
	q.mx.Unlock()
	go mqRun(q, handler)
	return true
}

func mqRun(q *MessageQueue, handler func(interface{}, bool)) {
	var curr *queueItem
	var stoped bool
	mx := q.mx
	cnd := q.cnd
	for {
		curr = nil
		//stoped = false
		mx.Lock()
		curr = q.head
		if curr == nil {
			cnd.Wait()
			curr = q.head
			if curr == nil {
				//if curr == nil stoped should be alway true. it means we woke up by signal from stop
				//and queue is empty
				mx.Unlock()
				if q.stopClbck != nil {
					q.stopClbck(q)
				}
				return
			}
		}
		q.cnt--
		stoped = q.state == mqsStoped
		q.head = curr.next
		mx.Unlock()

		handler(curr.data, !stoped)
		curr.free()
		if stoped {
			//handle all messages and quit
			mx.Lock()
			q.state = mqsStoped
			curr = q.head
			q.head = nil
			q.tail = nil
			q.cnt = 0
			mx.Unlock()
			for {
				if curr == nil {
					if q.stopClbck != nil {
						q.stopClbck(q)
					}
					return
				}
				handler(curr.data, false)
				curr.free()
				curr = curr.next
			}
		}
	}
}
