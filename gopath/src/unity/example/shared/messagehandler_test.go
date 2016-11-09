package shared

import "testing"

func TestMessageHandler(t *testing.T) {
	h := NewMessageHandler()
	m0 := newMessage(0, nil, true)
	m1 := newMessage(1, nil, true)
	m2 := newMessage(2, nil, true)
	var res int
	h.Handle(m2, true) // shouldn't raise
	if res != 0 {
		t.Error("wrong res with no handlers")
	}

	h.SetHandler(0, func(m *Message, isOn bool) {
		res = 0
	})

	h.SetHandler(1, func(m *Message, isOn bool) {
		res++
	})
	h.Handle(m1, true)
	if res == 0 {
		t.Error("code 1 was not handled")
	}
	h.Handle(m0, true)
	if res != 0 {
		t.Error("code 0 was not handled")
	}
	h.SetHandler(0, nil)
	h.Handle(m1, true)
	if res == 0 {
		t.Error("code 1 was not handled")
	}
	h.Handle(m0, true)
	if res == 0 {
		t.Error("code 0 shouldn't be handled after handler was delete")
	}

}
