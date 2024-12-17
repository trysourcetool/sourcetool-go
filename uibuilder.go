package sourcetool

import (
	"context"

	"github.com/trysourcetool/sourcetool-go/textinput"
)

type UIBuilder interface {
	Context() context.Context
	TextInput(string, ...textinput.Option) string
}

type uiBuilder struct {
	context context.Context
	cursor  *Cursor
	session *Session
	page    *Page
}

func (b *uiBuilder) Context() context.Context {
	return b.context
}
