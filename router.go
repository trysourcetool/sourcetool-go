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
	parent       *router
	sourcetool   *Sourcetool
	basePath     string
	namespaceDNS string
	groups       []string
}

func newRouter(st *Sourcetool, namespaceDNS string) Router {
	return &router{
		groups:       make([]string, 0),
		sourcetool:   st,
		namespaceDNS: namespaceDNS,
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

func (r *router) collectGroups() []string {
	groups := make([]string, 0)
	current := r
	for current != nil {
		groups = append(groups, current.groups...)
		current = current.parent
	}
	return groups
}

func (r *router) Page(path, name string, handler func(UIBuilder) error) {
	fullPath := r.joinPath(path)
	pageID := r.generatePageID(fullPath)

	page := &page{
		id:           pageID,
		name:         name,
		path:         fullPath,
		handler:      handler,
		accessGroups: removeDuplicates(r.collectGroups()),
	}

	r.sourcetool.addPage(pageID, page)
}

func (r *router) AccessGroups(groups ...string) Router {
	r.groups = append(r.groups, groups...)
	return r
}

func (r *router) Group(path string) Router {
	newRouter := &router{
		parent:       r,
		sourcetool:   r.sourcetool,
		basePath:     r.joinPath(path),
		namespaceDNS: r.namespaceDNS,
		groups:       make([]string, 0),
	}

	return newRouter
}
