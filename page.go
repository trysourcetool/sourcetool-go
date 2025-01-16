package sourcetool

import (
	"sync"

	"github.com/gofrs/uuid/v5"
)

type PageBuilder interface {
	Page(name string, handler func(UIBuilder) error) PageBuilder
	AccessGroups(groups ...string) PageBuilder
}

type pageBuilder struct {
	sourcetool   *Sourcetool
	page         *page
	namespaceDNS string
}

func (b *pageBuilder) currentPage() *page {
	return b.page
}

func (b *pageBuilder) Page(name string, handler func(UIBuilder) error) PageBuilder {
	b.page = &page{
		id:      b.generatePageID(name),
		name:    name,
		handler: handler,
	}

	b.addPage()

	return b
}

func (b *pageBuilder) AccessGroups(groups ...string) PageBuilder {
	b.page.accessGroups = groups
	return b
}

func (b *pageBuilder) generatePageID(pageName string) uuid.UUID {
	ns := uuid.NewV5(uuid.NamespaceDNS, b.namespaceDNS)
	return uuid.NewV5(ns, pageName)
}

func (b *pageBuilder) addPage() {
	b.sourcetool.mu.Lock()
	b.sourcetool.pages[b.page.id] = b.page
	b.sourcetool.mu.Unlock()
}

type page struct {
	id           uuid.UUID
	name         string
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
