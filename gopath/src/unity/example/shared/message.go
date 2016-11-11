package shared

import "sync"

//MessageSender desribes type that can sends messages
type MessageSender interface {
	Send(code MessageCode, data interface{})
	SendSync(code MessageCode, data interface{}) interface{}
}

//MessageCode specifies code of Message. Handlers could use to do some specific stuff
//but it's not required tu use it
type MessageCode byte

//Message is wrapper for message code and data. Messages could be joined in linked list
//it use internal pool, and returns automaticly to pool when handled, so don't create or cache it
//This code cause lead to very tricky results, because m returns to pool twice
//		m := New(Message)
//		h := NewMessageHandler()
//		h.Handle(m)
//		h.Handle(m)
type Message struct {
	//code of messsage. Provide additional info for handlers, but could be ignored
	Code MessageCode
	//additional data/ Could be nil
	Data interface{}
	//reference to next message in linked list
	//is message sent in sync mode
	sync bool
	//caller wait for this Cond if sync == true
	cnd *sync.Cond
	//reference to next Message if queued
	Next *Message
}

//message pool
var messageFree = sync.Pool{
	New: func() interface{} {
		return &Message{cnd: sync.NewCond(new(sync.Mutex))}
	},
}

// NewMessage allocates a new Message struct or grabs a cached one.
func NewMessage(code MessageCode, data interface{}, sync bool) *Message {
	m := messageFree.Get().(*Message)
	m.Code = code
	m.Data = data
	m.sync = sync
	if sync {
		m.cnd.L.Lock()
	}
	return m
}

//Free saves used pp structs in ppFree; avoids an allocation per invocation.
func (m *Message) free() {
	m.Next = nil
	m.Data = nil
	messageFree.Put(m)
}

//Handle calls hanler for Message and all messages queued through Next field
//if Message is sync - signals to waiter else frees message
func (m *Message) Handle(handler func(*Message)) {
	var tmp *Message
	for m != nil {
		handler(m)
		tmp = m
		m = m.Next
		if tmp.sync {
			tmp.cnd.L.Lock()
			tmp.cnd.Signal()
			tmp.cnd.L.Unlock()
		} else {
			tmp.free()
		}
	}
}

//Wait waits for sync message and returns data. Wait frees message
func (m *Message) Wait() interface{} {
	m.cnd.Wait()
	m.cnd.L.Unlock()
	data := m.Data
	m.free()
	return data
}
