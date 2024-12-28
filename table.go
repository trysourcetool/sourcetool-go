package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/table"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	tbl "github.com/trysourcetool/sourcetool-go/table"
)

func (b *uiBuilder) Table(data any, options ...table.Option) table.Value {
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
		return table.Value{}
	}
	page := b.page
	if page == nil {
		return table.Value{}
	}
	cursor := b.cursor
	if cursor == nil {
		return table.Value{}
	}
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateTableID(path)
	state := sess.State.GetTable(widgetID)
	if state == nil {
		state = &table.State{
			ID:    widgetID,
			Value: table.Value{},
		}
	}
	state.Data = data
	state.Header = opts.Header
	state.Description = opts.Description
	state.OnSelect = opts.OnSelect.String()
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: table.WidgetType,
		Path:       path,
		Data:       convertStateToTableData(state),
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateTableID(path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, table.WidgetType+"-"+path.String())
}

func convertStateToTableData(state *table.State) *websocket.TableData {
	if state == nil {
		return nil
	}
	data := &websocket.TableData{
		Data:         state.Data,
		Header:       state.Header,
		Description:  state.Description,
		OnSelect:     state.OnSelect,
		RowSelection: state.RowSelection,
		Value:        websocket.TableDataValue{},
	}
	if state.Value.Selection != nil {
		data.Value.Selection = &websocket.TableDataValueSelection{
			Row:  state.Value.Selection.Row,
			Rows: state.Value.Selection.Rows,
		}
	}
	return data
}

func convertTableDataToState(id uuid.UUID, data *websocket.TableData) *table.State {
	if data == nil {
		return nil
	}
	state := &table.State{
		ID:           id,
		Data:         data.Data,
		Header:       data.Header,
		Description:  data.Description,
		OnSelect:     data.OnSelect,
		RowSelection: data.RowSelection,
		Value:        table.Value{},
	}
	if data.Value.Selection != nil {
		state.Value.Selection = &table.Selection{
			Row:  data.Value.Selection.Row,
			Rows: data.Value.Selection.Rows,
		}
	}
	return state
}
