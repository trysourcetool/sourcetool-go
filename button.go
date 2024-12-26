package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/button"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

const widgetTypeButton = "button"

func (b *uiBuilder) Button(label string, options ...button.Option) bool {
	opts := &button.Options{
		Label:    label,
		Disabled: false,
	}

	for _, option := range options {
		option(opts)
	}

	sess := b.session
	if sess == nil {
		return false
	}
	page := b.page
	if page == nil {
		return false
	}
	cursor := b.cursor
	if cursor == nil {
		return false
	}
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateButtonInputID(label, path)
	state := sess.State.GetButton(widgetID)
	if state == nil {
		state = &button.State{
			ID:    widgetID,
			Value: false,
		}
	}
	state.Label = opts.Label
	state.Disabled = opts.Disabled
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeButton,
		Path:       path,
		Data:       state,
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateButtonInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, widgetTypeTextInput+"-"+label+"-"+path.String())
}
