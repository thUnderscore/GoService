package shared

import "sync"

//Switcher sync  bool flag and allows call func on flag changes
type Switcher struct {
	*sync.Mutex
	//Normaly use IsActive() function. Use this field if mutex is already locked
	Active bool
}

//On calls f if switched off and changes state of Switcher
func (s *Switcher) On(f func()) {
	s.Lock()
	defer s.Unlock()
	if s.Active {
		return
	}
	s.Active = true
	f()
}

//Off calls f if switched on and changes state of Switcher
func (s *Switcher) Off(f func()) {
	s.Lock()
	defer s.Unlock()
	if !s.Active {
		return
	}
	s.Active = false
	f()
}

//IsActive indicates if switcher is on
func (s *Switcher) IsActive() bool {
	s.Lock()
	defer s.Unlock()
	return s.Active
}

//NewSwitcher Creates new Switcher
func NewSwitcher() Switcher {
	return Switcher{Mutex: new(sync.Mutex)}
}
