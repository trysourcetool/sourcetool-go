package sourcetool

type navigation struct {
	name  string
	pages []*page
}

func (s *Sourcetool) Navigation(name string, handler func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentNav := &navigation{
		name: name,
	}
	s.navigations = append(s.navigations, currentNav)

	handler()
}
