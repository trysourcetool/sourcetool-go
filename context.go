package sourcetool

import (
	"context"
	"time"
)

type Context struct {
	context context.Context
	cursor  *Cursor
	session *Session
	page    *Page
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.context.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.context.Done()
}

func (c *Context) Err() error {
	return c.context.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.context.Value(key)
}
