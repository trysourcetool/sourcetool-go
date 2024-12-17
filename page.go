package sourcetool

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type Page struct {
	ID      uuid.UUID
	Name    string
	Handler func(*Context) error
	Context context.Context
}

func (p *Page) Run(ctx *Context) error {
	ctx.context = p.Context
	if err := p.Handler(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Sourcetool) Page(ctx context.Context, name string, handler func(*Context) error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := &Page{
		ID:      s.generatePageID(name),
		Name:    name,
		Handler: handler,
		Context: ctx,
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
