package sourcetool

import (
	"context"

	"github.com/trysourcetool/sourcetool-go/button"
	"github.com/trysourcetool/sourcetool-go/table"
	"github.com/trysourcetool/sourcetool-go/textinput"
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
	cursor  *Cursor
	session *Session
	page    *Page
}

func (b *uiBuilder) Context() context.Context {
	return b.context
}
