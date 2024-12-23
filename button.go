package sourcetool

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/trysourcetool/sourcetool-go/button"
	ws "github.com/trysourcetool/sourcetool-go/websocket"
)

const widgetTypeButton = "button"

func (b *uiBuilder) Button(label string, options ...button.Option) bool {
	opts := &button.Options{
		Label: label,
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
	path := cursor.GetDeltaPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.ID.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateButtonInputID(label, path)

	log.Printf("Button ID: %s\n", widgetID.String())

	state := sess.State.GetButton(widgetID)
	if state == nil {
		// Set initial state
		state = &button.State{
			ID:       widgetID,
			Label:    opts.Label,
			Disabled: opts.Disabled,
		}
		sess.State.Set(widgetID, state)
	}
	returnValue := state.Value

	b.runtime.EnqueueMessage(uuid.Must(uuid.NewV4()).String(), ws.MessageMethodRenderWidget, &ws.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.ID.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeButton,
		Data:       state,
	})

	cursor.Next()

	return bool(returnValue)
}

func (b *uiBuilder) generateButtonInputID(label string, path []int) uuid.UUID {
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
