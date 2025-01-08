package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/table"
)

func (b *uiBuilder) Table(data any, opts ...table.Option) table.Value {
	tableOpts := &options.TableOptions{
		OnSelect:     conv.NilValue(table.SelectionBehaviorIgnore.String()),
		RowSelection: conv.NilValue(table.SelectionModeSingle.String()),
	}

	for _, o := range opts {
		o.Apply(tableOpts)
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
	tableState := sess.State.GetTable(widgetID)
	if tableState == nil {
		tableState = &state.TableState{
			ID:    widgetID,
			Value: state.TableStateValue{},
		}
	}
	tableState.Data = data
	tableState.Header = tableOpts.Header
	tableState.Description = tableOpts.Description
	tableState.OnSelect = tableOpts.OnSelect
	tableState.RowSelection = tableOpts.RowSelection
	sess.State.Set(widgetID, tableState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeTable.String(),
		Path:       path,
		Data:       convertStateToTableData(tableState),
	})

	cursor.next()

	value := table.Value{}
	if tableState.Value.Selection != nil {
		value.Selection = &table.Selection{
			Row:  tableState.Value.Selection.Row,
			Rows: tableState.Value.Selection.Rows,
		}
	}

	return value
}

func (b *uiBuilder) generateTableID(path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeTable.String()+"-"+path.String())
}

func convertStateToTableData(state *state.TableState) *websocket.TableData {
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

func convertTableDataToState(id uuid.UUID, data *websocket.TableData) *state.TableState {
	if data == nil {
		return nil
	}
	tableState := &state.TableState{
		ID:           id,
		Data:         data.Data,
		Header:       data.Header,
		Description:  data.Description,
		OnSelect:     data.OnSelect,
		RowSelection: data.RowSelection,
		Value:        state.TableStateValue{},
	}
	if data.Value.Selection != nil {
		tableState.Value.Selection = &state.TableStateValueSelection{
			Row:  data.Value.Selection.Row,
			Rows: data.Value.Selection.Rows,
		}
	}
	return tableState
}
