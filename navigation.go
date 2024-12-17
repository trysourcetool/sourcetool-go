package sourcetool

type Navigation struct {
	Name  string
	Pages []*Page
}

func (s *Sourcetool) Navigation(name string, handler func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentNav := &Navigation{
		Name: name,
	}
	s.navigations = append(s.navigations, currentNav)

	handler()
}
