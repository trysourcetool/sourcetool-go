package sourcetool

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/textinput"
	ws "github.com/trysourcetool/sourcetool-go/websocket"
)

const widgetTypeTextInput = "textInput"

func (b *uiBuilder) TextInput(label string, options ...textinput.Option) string {
	opts := &textinput.Options{
		Label: label,
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
	path := cursor.GetDeltaPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.ID.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateTextInputID(label, path)

	log.Printf("Text Input ID: %s\n", widgetID.String())

	state := sess.State.GetTextInput(widgetID)
	if state == nil {
		// Set initial state
		state = &textinput.State{
			ID:           widgetID,
			Label:        opts.Label,
			Value:        textinput.ReturnValue(opts.DefaultValue),
			Placeholder:  opts.Placeholder,
			DefaultValue: opts.DefaultValue,
			Required:     opts.Required,
			MaxLength:    opts.MaxLength,
			MinLength:    opts.MinLength,
		}
		sess.State.Set(widgetID, state)
	}
	returnValue := state.Value

	b.runtime.EnqueueMessage(uuid.Must(uuid.NewV4()).String(), ws.MessageMethodRenderWidget, &ws.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.ID.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeTextInput,
		Data:       state,
	})

	cursor.Next()

	return string(returnValue)
}

func (b *uiBuilder) generateTextInputID(label string, path []int) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	strPath := make([]string, len(path))
	for i, num := range path {
		strPath[i] = fmt.Sprint(num)
	}
	return uuid.NewV5(page.ID, widgetTypeTextInput+"-"+label+"-"+strings.Join(strPath, ""))
}
