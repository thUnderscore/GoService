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
//		m := new(Message)
//		h := NewMessageHandler()
//		h.Handle(m)
//		h.Handle(m)
type Message struct {
	//code of messsage. Provide additional info for handlers, but could be ignored
	code MessageCode
	//is message sent in sync mode
	sync bool
	//caller wait for this Cond if sync == true
	cnd *sync.Cond
	//additional data/ Could be nil
	data interface{}
	//reference to next message in linked list
	next *Message
}

//message pool
var messageFree = sync.Pool{
	New: func() interface{} {
		return &Message{cnd: sync.NewCond(new(sync.Mutex))}
	},
}

// newPrinter allocates a new pp struct or grabs a cached one.
func newMessage(code MessageCode, data interface{}, sync bool) *Message {
	m := messageFree.Get().(*Message)
	m.code = code
	m.data = data
	m.sync = sync
	if sync {
		m.cnd.L.Lock()
	}
	return m
}

// free saves used pp structs in ppFree; avoids an allocation per invocation.
func (m *Message)free() {
	m.next = nil
	m.data = nil
	messageFree.Put(m)
}

func (m *Message)handle(handler func(*Message)) {
	var tmp *Message
	for m != nil {
		handler(m)
		tmp = m
		m = m.next
		if tmp.sync {
			tmp.cnd.L.Lock()
			tmp.cnd.Signal()
			tmp.cnd.L.Unlock()
		} else {
			tmp.free()
		}
	}
}

func (m *Message) wait() interface{} {
	m.cnd.Wait()
	m.cnd.L.Unlock()
	data := m.data
	m.free()
	return data
}
