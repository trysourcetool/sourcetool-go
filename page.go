package sourcetool

import (
	"fmt"
	"sync"

	"github.com/gofrs/uuid/v5"
)

type page struct {
	id      uuid.UUID
	name    string
	handler func(UIBuilder) error
}

func (p *page) run(ui UIBuilder) error {
	if err := p.handler(ui); err != nil {
		return err
	}
	return nil
}

func (s *Sourcetool) Page(name string, handler func(UIBuilder) error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := &page{
		id:      s.generatePageID(name),
		name:    name,
		handler: handler,
	}

	s.pages[p.id] = p
	if len(s.navigations) > 0 {
		currentNav := s.navigations[len(s.navigations)-1]
		currentNav.pages = append(currentNav.pages, p)
	}
}

func (s *Sourcetool) generatePageID(pageName string) uuid.UUID {
	ns := uuid.NewV5(uuid.NamespaceDNS, fmt.Sprintf("%s.trysourcetool.com", s.subdomain))
	return uuid.NewV5(ns, pageName)
}

type pageManager struct {
	pages map[uuid.UUID]*page
	mu    sync.RWMutex
}

func newPageManager(pages map[uuid.UUID]*page) *pageManager {
	return &pageManager{
		pages: pages,
	}
}

func (s *pageManager) getPage(id uuid.UUID) *page {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.pages[id]
}
