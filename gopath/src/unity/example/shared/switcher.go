package shared

import "sync"

//Switcher sync  bool flag and allows call func on flag changes
type Switcher struct {
	mx *sync.Mutex
	on bool
}

//On calls f if switched off and changes state of Switcher
func (s *Switcher) On(f func()) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.on {
		return
	}
	s.on = true
	f()
}

//Off calls f if switched on and changes state of Switcher
func (s *Switcher) Off(f func()) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if !s.on {
		return
	}
	s.on = false
	f()
}

func (s *Switcher) isOn() bool {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.on
}

//NewSwitcher Creates new Switcher
func NewSwitcher() Switcher {
	return Switcher{mx: new(sync.Mutex)}
}
