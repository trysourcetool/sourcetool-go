package sourcetool

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/table"
	ws "github.com/trysourcetool/sourcetool-go/websocket"
)

func (b *uiBuilder) Table(data any, options ...table.Option) table.ReturnValue {
	opts := &table.Options{}

	for _, option := range options {
		option(opts)
	}

	sess := b.session
	if sess == nil {
		return table.ReturnValue{}
	}
	page := b.page
	if page == nil {
		return table.ReturnValue{}
	}
	cursor := b.cursor
	if cursor == nil {
		return table.ReturnValue{}
	}
	path := cursor.GetDeltaPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.ID.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateTableID(path)

	log.Printf("Table ID: %s\n", widgetID.String())

	var returnValue table.ReturnValue
	state := sess.State.GetTable(widgetID)
	if state == nil {
		// Set initial state
		state = &table.State{
			ID:          widgetID,
			Data:        data,
			Header:      opts.Header,
			Description: opts.Description,
			OnSelect:    opts.OnSelect.String(),
		}
		sess.State.Set(widgetID, state)
	} else {
		returnValue = state.Value
	}

	b.runtime.EnqueueMessage(uuid.Must(uuid.NewV4()).String(), ws.MessageMethodRenderWidget, &ws.RenderWidgetPayload{
		SessionID: sess.ID.String(),
		PageID:    page.ID.String(),
		WidgetID:  widgetID.String(),
		Data:      state,
	})

	cursor.Next()

	return returnValue
}

func (b *uiBuilder) generateTableID(path []int) uuid.UUID {
	const widgetType = "table"
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	strPath := make([]string, len(path))
	for i, num := range path {
		strPath[i] = fmt.Sprint(num)
	}
	return uuid.NewV5(page.ID, widgetType+"-"+strings.Join(strPath, ""))
}
