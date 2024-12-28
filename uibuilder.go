package sourcetool

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/trysourcetool/sourcetool-go/internal/button"
	"github.com/trysourcetool/sourcetool-go/internal/columns"
	"github.com/trysourcetool/sourcetool-go/internal/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/form"
	"github.com/trysourcetool/sourcetool-go/internal/numberinput"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/table"
	"github.com/trysourcetool/sourcetool-go/internal/textarea"
	"github.com/trysourcetool/sourcetool-go/internal/textinput"
)

type UIBuilder interface {
	Context() context.Context
	TextInput(string, ...textinput.Option) string
	NumberInput(string, ...numberinput.Option) *float64
	DateInput(string, ...dateinput.Option) *time.Time
	TextArea(string, ...textarea.Option) string
	Table(any, ...table.Option) table.Value
	Button(string, ...button.Option) bool
	Form(string, ...form.Option) (UIBuilder, bool)
	Columns(int, ...columns.Option) []UIBuilder
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

type path []int

func (p path) String() string {
	strPath := make([]string, len(p))
	for i, num := range p {
		strPath[i] = fmt.Sprint(num)
	}
	return strings.Join(strPath, "")
}

type cursor struct {
	parentPath []int
	index      int
}

func newCursor() *cursor {
	return &cursor{
		parentPath: []int{},
		index:      0,
	}
}

func (c *cursor) getPath() path {
	path := make([]int, len(c.parentPath))
	copy(path, c.parentPath)
	path = append(path, c.index)
	return path
}

func (c *cursor) next() {
	c.index++
}
