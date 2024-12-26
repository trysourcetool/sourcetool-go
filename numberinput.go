package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/numberinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

const widgetTypeNumberInput = "numberInput"

func (b *uiBuilder) NumberInput(label string, options ...numberinput.Option) float64 {
	opts := &numberinput.Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: 0,
		Required:     false,
		MaxValue:     nil,
		MinValue:     nil,
	}

	for _, option := range options {
		option(opts)
	}

	sess := b.session
	if sess == nil {
		return 0
	}
	page := b.page
	if page == nil {
		return 0
	}
	cursor := b.cursor
	if cursor == nil {
		return 0
	}
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateNumberInputID(label, path)
	state := sess.State.GetNumberInput(widgetID)
	if state == nil {
		state = &numberinput.State{
			ID:    widgetID,
			Value: opts.DefaultValue,
		}
	}
	state.Label = opts.Label
	state.Placeholder = opts.Placeholder
	state.DefaultValue = opts.DefaultValue
	state.Required = opts.Required
	state.MaxValue = opts.MaxValue
	state.MinValue = opts.MinValue
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeNumberInput,
		Path:       path,
		Data:       state,
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateNumberInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, widgetTypeNumberInput+"-"+label+"-"+path.String())
}
