package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/textinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

const widgetTypeTextInput = "textInput"

func (b *uiBuilder) TextInput(label string, options ...textinput.Option) string {
	opts := &textinput.Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: "",
		Required:     false,
		MaxLength:    nil,
		MinLength:    nil,
	}

	for _, option := range options {
		option(opts)
	}

	sess := b.session
	if sess == nil {
		return ""
	}
	page := b.page
	if page == nil {
		return ""
	}
	cursor := b.cursor
	if cursor == nil {
		return ""
	}
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateTextInputID(label, path)
	state := sess.State.GetTextInput(widgetID)
	if state == nil {
		state = &textinput.State{
			ID:    widgetID,
			Value: textinput.ReturnValue(opts.DefaultValue),
		}
	}
	state.Label = opts.Label
	state.Placeholder = opts.Placeholder
	state.DefaultValue = opts.DefaultValue
	state.Required = opts.Required
	state.MaxLength = opts.MaxLength
	state.MinLength = opts.MinLength
	sess.State.Set(widgetID, state)

	returnValue := state.Value

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeTextInput,
		Path:       path,
		Data:       state,
	})

	cursor.next()

	return string(returnValue)
}

func (b *uiBuilder) generateTextInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, widgetTypeTextInput+"-"+label+"-"+path.String())
}
