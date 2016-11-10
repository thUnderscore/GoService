package shared

import "testing"

func TestMessageHandler(t *testing.T) {
	h := NewMessageHandler()
	
	var res int
	h.Handle(newMessage(2, nil, false)) // shouldn't raise
	if res != 0 {
		t.Error("wrong res with no handlers")
	}

	h.SetHandler(0, func(m *Message) {
		res = 0
	})

	h.SetHandler(1, func(m *Message) {
		res++
	})
	h.Handle(newMessage(1, nil, false))
	if res == 0 {
		t.Error("code 1 was not handled")
	}
	h.Handle(newMessage(0, nil, false))
	if res != 0 {
		t.Error("code 0 was not handled")
	}
	h.SetHandler(0, nil)
	h.Handle(newMessage(1, nil, false))
	if res == 0 {
		t.Error("code 1 was not handled")
	}
	h.Handle(newMessage(0, nil, false))
	if res == 0 {
		t.Error("code 0 shouldn't be handled after handler was delete")
	}

}
