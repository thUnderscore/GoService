package shared

import "sync"

type MessageCode int

type Message struct {
	code MessageCode
	sync bool
	cnd  *sync.Cond
	data interface{}
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
func (m *Message) free() {
	m.next = nil
	m.data = nil
	messageFree.Put(m)
}

func (m *Message) handled() {
	if m.sync {
		m.cnd.L.Lock()
		m.cnd.Signal()
		m.cnd.L.Unlock()
	} else {
		m.free()
	}
}

func (m *Message) wait() interface{} {
	m.cnd.Wait()
	m.cnd.L.Unlock()
	data := m.data
	m.free()
	return data
}
