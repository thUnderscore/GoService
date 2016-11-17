package shared

import "sync"

//Switcher sync  bool flag and allows call func on flag changes
type Switcher struct {
	*sync.Mutex
	//Normaly use IsActive() function. Use this field if mutex is already locked
	Active bool
	//Indicate if was switched on at least one time already
	Used bool
}

//On calls f if switched off and changes state of Switcher
//Doesn't do it if once == true and was switched on at least one time already
func (s *Switcher) On(f func(), once bool) {
	s.Lock()
	defer s.Unlock()
	if s.Active || (s.Used && once) {
		return
	}
	s.Used = true
	s.Active = true
	if f != nil {
		f()
	}
}

//Off calls f if switched on and changes state of Switcher
func (s *Switcher) Off(f func()) {
	s.Lock()
	defer s.Unlock()
	if !s.Active {
		return
	}
	s.Active = false
	if f != nil {
		f()
	}
}

//IsActive indicates if switcher is on
func (s *Switcher) IsActive() bool {
	s.Lock()
	defer s.Unlock()
	return s.Active
}

//NewSwitcher Creates new Switcher
func NewSwitcher() *Switcher {
	return &Switcher{Mutex: new(sync.Mutex)}
}
