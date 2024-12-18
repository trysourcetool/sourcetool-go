package sourcetool

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/textinput"
	ws "github.com/trysourcetool/sourcetool-go/websocket"
)

func (ui *uiBuilder) TextInput(label string, options ...textinput.Option) string {
	opts := &textinput.Options{
		Label: label,
	}

	for _, option := range options {
		option(opts)
	}

	sess := ui.session
	if sess == nil {
		return ""
	}
	page := ui.page
	if page == nil {
		return ""
	}
	cursor := ui.cursor
	if cursor == nil {
		return ""
	}
	path := cursor.GetDeltaPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.ID.String())
	log.Printf("Path: %v\n", path)

	id := generateTextInputID(page.ID, label, path)

	log.Printf("Text Input ID: %s\n", id.String())

	var returnValue string
	state := sess.State.GetTextInput(id)
	if state == nil {
		// Set initial state
		state = &textinput.State{
			ID:           id,
			Label:        opts.Label,
			Placeholder:  opts.Placeholder,
			DefaultValue: opts.DefaultValue,
			Required:     opts.Required,
			MaxLength:    opts.MaxLength,
			MinLength:    opts.MinLength,
		}
		sess.State.Set(id, state)
	} else {
		returnValue = state.Value
	}

	if Runtime == nil {
		return ""
	}

	Runtime.EnqueueMessage(uuid.Must(uuid.NewV4()).String(), ws.MessageMethodRenderWidget, &ws.RenderWidgetPayload{
		SessionID: sess.ID.String(),
		PageID:    page.ID.String(),
		WidgetID:  id.String(),
		Data:      state,
	})

	cursor.Next()

	return returnValue
}

func generateTextInputID(pageID uuid.UUID, label string, path []int) uuid.UUID {
	strPath := make([]string, len(path))
	for i, num := range path {
		strPath[i] = fmt.Sprint(num)
	}
	return uuid.NewV5(pageID, label+"-"+strings.Join(strPath, ""))
}
