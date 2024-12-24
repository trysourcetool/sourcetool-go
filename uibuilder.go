package sourcetool

import (
	"context"

	"github.com/trysourcetool/sourcetool-go/internal/button"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/table"
	"github.com/trysourcetool/sourcetool-go/internal/textinput"
)

type UIBuilder interface {
	Context() context.Context
	TextInput(string, ...textinput.Option) string
	Table(any, ...table.Option) table.ReturnValue
	Button(string, ...button.Option) bool
}

type uiBuilder struct {
	runtime *runtime
	context context.Context
	cursor  *cursor
	session *session.Session
	page    *page
}

func (b *uiBuilder) Context() context.Context {
	return b.context
}

type containerType int

const main containerType = iota

type cursor struct {
	rootContainer containerType
	parentPath    []int
	index         int
}

func newCursor(container containerType) *cursor {
	return &cursor{
		rootContainer: container,
		parentPath:    []int{},
		index:         0,
	}
}

func (c *cursor) getDeltaPath() []int {
	path := []int{int(c.rootContainer)}
	path = append(path, c.parentPath...)
	path = append(path, c.index)
	return path
}

func (c *cursor) next() {
	c.index++
}
