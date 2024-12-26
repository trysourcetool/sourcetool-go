package sourcetool

import (
	"github.com/gofrs/uuid/v5"
	"github.com/trysourcetool/sourcetool-go/internal/form"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

const widgetTypeForm = "form"

func (b *uiBuilder) Form(options ...form.Option) UIBuilder {
	opts := &form.Options{
		ClearOnSubmit: false,
	}

	for _, option := range options {
		option(opts)
	}

	sess := b.session
	if sess == nil {
		return b
	}
	page := b.page
	if page == nil {
		return b
	}
	cursor := b.cursor
	if cursor == nil {
		return b
	}
	path := cursor.getPath()

	widgetID := b.generateFormID(path)
	state := sess.State.GetForm(widgetID)
	if state == nil {
		state = &form.State{
			ID:            widgetID,
			ClearOnSubmit: opts.ClearOnSubmit,
		}
		sess.State.Set(widgetID, state)
	}

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeForm,
		Path:       path,
		Data:       state,
	})

	cursor.next()

	childCursor := newCursor()
	childCursor.parentPath = path

	childBuilder := &uiBuilder{
		runtime: b.runtime,
		session: sess,
		page:    page,
		cursor:  childCursor,
	}

	return childBuilder
}

func (b *uiBuilder) generateFormID(path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, widgetTypeForm+"-"+path.String())
}
