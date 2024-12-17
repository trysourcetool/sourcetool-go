package sourcetool

import (
	"sync"

	"github.com/gofrs/uuid/v5"
)

type PageManager struct {
	pages map[uuid.UUID]*Page
	mu    sync.RWMutex
}

func NewPageManager(pages map[uuid.UUID]*Page) *PageManager {
	return &PageManager{
		pages: pages,
	}
}

func (s *PageManager) SetPage(page *Page) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pages[page.ID] = page
}

func (s *PageManager) GetPage(id uuid.UUID) *Page {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.pages[id]
}
