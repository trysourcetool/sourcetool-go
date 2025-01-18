package sourcetool

import (
	"sync"

	"github.com/gofrs/uuid/v5"
)

type page struct {
	id           uuid.UUID
	name         string
	path         string
	handler      func(UIBuilder) error
	accessGroups []string
}

func (p *page) run(ui UIBuilder) error {
	if err := p.handler(ui); err != nil {
		return err
	}
	return nil
}

func (p *page) hasAccess(userGroups []string) bool {
	if len(p.accessGroups) == 0 {
		return true
	}

	for _, userGroup := range userGroups {
		for _, requiredGroup := range p.accessGroups {
			if userGroup == requiredGroup {
				return true
			}
		}
	}
	return false
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
