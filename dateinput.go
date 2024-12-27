package sourcetool

import (
	"log"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) DateInput(label string, options ...dateinput.Option) *time.Time {
	opts := &dateinput.Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Format:       "YYYY/MM/DD",
		MaxValue:     nil,
		MinValue:     nil,
		Location:     time.Local,
	}

	for _, option := range options {
		option(opts)
	}

	sess := b.session
	if sess == nil {
		return nil
	}
	page := b.page
	if page == nil {
		return nil
	}
	cursor := b.cursor
	if cursor == nil {
		return nil
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
			Value: opts.DefaultValue,
		}
	}
	state.Label = opts.Label
	state.Placeholder = opts.Placeholder
	state.DefaultValue = opts.DefaultValue
	state.Required = opts.Required
	state.Format = opts.Format
	state.MaxValue = opts.MaxValue
	state.MinValue = opts.MinValue
	state.Location = opts.Location
	sess.State.Set(widgetID, state)

	var value, defaultValue, maxValue, minValue string
	if state.Value != nil {
		value = state.Value.Format(time.DateOnly)
	}
	if state.DefaultValue != nil {
		defaultValue = state.DefaultValue.Format(time.DateOnly)
	}
	if state.MaxValue != nil {
		maxValue = state.MaxValue.Format(time.DateOnly)
	}
	if state.MinValue != nil {
		minValue = state.MinValue.Format(time.DateOnly)
	}
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: dateinput.WidgetType,
		Path:       path,
		Data: &websocket.DateInputData{
			Value:        value,
			Label:        state.Label,
			Placeholder:  state.Placeholder,
			DefaultValue: defaultValue,
			Required:     state.Required,
			Format:       state.Format,
			MaxValue:     maxValue,
			MinValue:     minValue,
		},
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateDateInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, dateinput.WidgetType+"-"+label+"-"+path.String())
}
