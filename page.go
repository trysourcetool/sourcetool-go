package sourcetool

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type Page struct {
	ID      uuid.UUID
	Name    string
	Handler func(UIBuilder) error
}

func (p *Page) Run(ui UIBuilder) error {
	if err := p.Handler(ui); err != nil {
		return err
	}
	return nil
}

func (s *Sourcetool) Page(name string, handler func(UIBuilder) error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := &Page{
		ID:      s.generatePageID(name),
		Name:    name,
		Handler: handler,
	}

	s.pages[p.ID] = p
	if len(s.navigations) > 0 {
		currentNav := s.navigations[len(s.navigations)-1]
		currentNav.Pages = append(currentNav.Pages, p)
	}
}

func (s *Sourcetool) generatePageID(pageName string) uuid.UUID {
	ns := uuid.NewV5(uuid.NamespaceDNS, fmt.Sprintf("%s.trysourcetool.com", s.subdomain))
	return uuid.NewV5(ns, pageName)
}
