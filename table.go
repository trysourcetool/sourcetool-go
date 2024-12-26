package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/table"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	tbl "github.com/trysourcetool/sourcetool-go/table"
)

const widgetTypeTable = "table"

func (b *uiBuilder) Table(data any, options ...table.Option) table.ReturnValue {
	opts := &table.Options{
		Header:       "",
		Description:  "",
		OnSelect:     tbl.OnSelectIgnore,
		RowSelection: tbl.RowSelectionSingle,
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
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateTableID(path)
	state := sess.State.GetTable(widgetID)
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
		sess.State.Set(widgetID, state)
	}
	state.Data = data
	returnValue := state.Value

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeTable,
		Path:       path,
		Data:       state,
	})

	cursor.next()

	return returnValue
}

func (b *uiBuilder) generateTableID(path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, widgetTypeTable+"-"+path.String())
}
