package sourcetool

import (
	"strings"

	"github.com/gofrs/uuid/v5"
)

type Router interface {
	Page(path, name string, handler func(UIBuilder) error)
	AccessGroups(groups ...string) Router
	Group(path string) Router
}

type router struct {
	parent         *router
	sourcetool     *Sourcetool
	basePath       string
	namespaceDNS   string
	routerGroups   []string
	pageGroups     []string
	nextIsPageCall bool
	lastWasGroup   bool
}

func newRouter(st *Sourcetool, namespaceDNS string) Router {
	return &router{
		routerGroups:   make([]string, 0),
		pageGroups:     make([]string, 0),
		sourcetool:     st,
		namespaceDNS:   namespaceDNS,
		nextIsPageCall: false,
		lastWasGroup:   false,
	}
}

func (r *router) generatePageID(fullPath string) uuid.UUID {
	ns := uuid.NewV5(uuid.NamespaceDNS, r.namespaceDNS)
	return uuid.NewV5(ns, fullPath)
}

func (r *router) joinPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if r.basePath == "" {
		return path
	}
	basePath := strings.TrimSuffix(r.basePath, "/")
	cleanPath := strings.TrimPrefix(path, "/")
	return basePath + "/" + cleanPath
}

func removeDuplicates(groups []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(groups))

	for _, group := range groups {
		if _, exists := seen[group]; !exists {
			seen[group] = struct{}{}
			result = append(result, group)
		}
	}
	return result
}

func (r *router) collectRouterGroups() []string {
	groups := make([]string, 0)
	current := r
	for current != nil {
		groups = append(groups, current.routerGroups...)
		current = current.parent
	}
	return groups
}

func (r *router) Page(path, name string, handler func(UIBuilder) error) {
	fullPath := r.joinPath(path)
	pageID := r.generatePageID(fullPath)

	pageGroups := make([]string, 0)

	pageGroups = append(pageGroups, r.collectRouterGroups()...)

	if r.nextIsPageCall && !r.lastWasGroup {
		pageGroups = append(pageGroups, r.pageGroups...)
	}

	r.nextIsPageCall = false
	r.lastWasGroup = false
	r.pageGroups = make([]string, 0)

	page := &page{
		id:           pageID,
		name:         name,
		path:         fullPath,
		handler:      handler,
		accessGroups: removeDuplicates(pageGroups),
	}

	r.sourcetool.addPage(pageID, page)
}

func (r *router) AccessGroups(groups ...string) Router {
	if r.lastWasGroup {
		r.routerGroups = append(r.routerGroups, groups...)
		r.nextIsPageCall = false
	} else {
		r.pageGroups = append(r.pageGroups, groups...)
		r.nextIsPageCall = true
	}
	r.lastWasGroup = false
	return r
}

func (r *router) Group(path string) Router {
	newRouter := &router{
		parent:         r,
		sourcetool:     r.sourcetool,
		basePath:       r.joinPath(path),
		routerGroups:   make([]string, 0),
		pageGroups:     make([]string, 0),
		namespaceDNS:   r.namespaceDNS,
		nextIsPageCall: false,
		lastWasGroup:   true,
	}
	return newRouter
}
