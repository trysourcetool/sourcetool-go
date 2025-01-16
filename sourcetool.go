package sourcetool

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gofrs/uuid/v5"
)

type Sourcetool struct {
	apiKey    string
	endpoint  string
	subdomain string
	runtime   *runtime
	pages     map[uuid.UUID]*page
	mu        sync.RWMutex
}

func New(apiKey string) *Sourcetool {
	subdomain := strings.Split(apiKey, "_")[0]
	s := &Sourcetool{
		apiKey:    apiKey,
		subdomain: subdomain,
		endpoint:  fmt.Sprintf("ws://%s.local.trysourcetool.com:8080/ws", subdomain),
		pages:     make(map[uuid.UUID]*page),
	}
	return s
}

func (s *Sourcetool) Listen() error {
	r, err := startRuntime(s.apiKey, s.endpoint, s.pages)
	if err != nil {
		return err
	}
	defer r.wsClient.Close()

	s.runtime = r

	return r.wsClient.Wait()
}

func (s *Sourcetool) Close() error {
	return s.runtime.wsClient.Close()
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

func (p *page) AccessGroups(groups ...string) *page {
	p.accessGroups = groups
	return p
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

func (s *Sourcetool) Page(name string, handler func(UIBuilder) error) *page {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := &page{
		id:      s.generatePageID(name),
		name:    name,
		handler: handler,
	}

	s.pages[p.id] = p

	return p
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
