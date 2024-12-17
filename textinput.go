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

	if Runtime == nil {
		return ""
	}

	var returnValue string
	if val := sess.State.Get(id); val != nil {
		if opt, ok := val.(*textinput.State); ok {
			returnValue = opt.Value
		}
	}

	Runtime.EnqueueMessage(uuid.Must(uuid.NewV4()).String(), ws.MessageMethodRenderWidget, &ws.RenderWidgetPayload{
		WidgetID:  id.String(),
		SessionID: sess.ID.String(),
		PageID:    page.ID.String(),
		Data:      opts,
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
