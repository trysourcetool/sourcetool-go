package sourcetool

import (
	"log"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

const widgetTypeDateInput = "dateInput"

func (b *uiBuilder) DateInput(label string, options ...dateinput.Option) time.Time {
	opts := &dateinput.Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Format:       "YYYY/MM/DD",
		MaxValue:     nil,
		MinValue:     nil,
	}

	for _, option := range options {
		option(opts)
	}

	sess := b.session
	if sess == nil {
		return time.Time{}
	}
	page := b.page
	if page == nil {
		return time.Time{}
	}
	cursor := b.cursor
	if cursor == nil {
		return time.Time{}
	}
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateDateInputID(label, path)
	state := sess.State.GetDateInput(widgetID)
	if state == nil {
		state = &dateinput.State{
			ID:    widgetID,
			Value: conv.SafeValue(opts.DefaultValue),
		}
	}
	state.Label = opts.Label
	state.Placeholder = opts.Placeholder
	state.DefaultValue = opts.DefaultValue
	state.Required = opts.Required
	state.Format = opts.Format
	state.MaxValue = opts.MaxValue
	state.MinValue = opts.MinValue
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeDateInput,
		Path:       path,
		Data:       state,
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateDateInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, widgetTypeDateInput+"-"+label+"-"+path.String())
}
