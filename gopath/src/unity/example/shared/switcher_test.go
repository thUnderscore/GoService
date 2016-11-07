package shared

import "testing"

func TestSwitcher(t *testing.T) {
	s := NewSwitcher()
	c := 0
	for i := 0; i < 20; i++ {
		s.On(func() {
			c++
		})
	}
	Sleep100ms()
	if c != 1 {
		t.Error("counter should be equal", 1)
	}
	if !s.isOn() {
		t.Error("switcher should be on")
	}
	for i := 0; i < 20; i++ {
		s.Off(func() {
			c = c + 10
		})
	}
	Sleep100ms()
	if s.isOn() {
		t.Error("switcher should be off")
	}
	if c != 11 {
		t.Error("counter should be equal", 11)
	}
	for i := 0; i < 20; i++ {
		s.On(func() {
			c++
		})
	}
	Sleep100ms()
	if c != 12 {
		t.Error("counter should be equal", 12)
	}
	if !s.isOn() {
		t.Error("switcher should be on")
	}

}
