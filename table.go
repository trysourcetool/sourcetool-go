package sourcetool

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/table"
)

const widgetTypeTable = "table"

func (b *uiBuilder) Table(data any, options ...table.Option) table.ReturnValue {
	opts := &table.Options{
		Header:       "",
		Description:  "",
		OnSelect:     table.OnSelectIgnore,
		RowSelection: table.RowSelectionSingle,
	}

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
	path := cursor.getDeltaPath()

	log.Printf("Session ID: %s", sess.id.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateTableID(path)

	log.Printf("Table ID: %s\n", widgetID.String())

	state := sess.state.getTable(widgetID)
	if state == nil {
		// Set initial state
		state = &table.State{
			ID:          widgetID,
			Value:       table.ReturnValue{},
			Data:        data,
			Header:      opts.Header,
			Description: opts.Description,
			OnSelect:    opts.OnSelect.String(),
		}
		sess.state.set(widgetID, state)
	}
	returnValue := state.Value

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.id.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeTable,
		Data:       state,
	})

	cursor.next()

	return returnValue
}

func (b *uiBuilder) generateTableID(path []int) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	strPath := make([]string, len(path))
	for i, num := range path {
		strPath[i] = fmt.Sprint(num)
	}
	return uuid.NewV5(page.id, widgetTypeTable+"-"+strings.Join(strPath, ""))
}
