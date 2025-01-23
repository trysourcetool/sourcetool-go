package sourcetool

import (
	"encoding/json"
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/table"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
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

	tableProto, err := convertStateToTableProto(tableState)
	if err != nil {
		return table.Value{}
	}
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_Table{
				Table: tableProto,
			},
		},
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

func convertStateToTableProto(state *state.TableState) (*widgetv1.Table, error) {
	if state == nil {
		return nil, nil
	}
	dataBytes, err := json.Marshal(state.Data)
	if err != nil {
		return nil, err
	}
	data := &widgetv1.Table{
		Data:         dataBytes,
		Header:       state.Header,
		Description:  state.Description,
		OnSelect:     state.OnSelect,
		RowSelection: state.RowSelection,
		Value:        &widgetv1.TableValue{},
	}
	if state.Value.Selection != nil {
		rows := make([]uint32, len(state.Value.Selection.Rows))
		for i, r := range state.Value.Selection.Rows {
			rows[i] = uint32(r)
		}
		data.Value.Selection = &widgetv1.TableValueSelection{
			Row:  uint32(state.Value.Selection.Row),
			Rows: rows,
		}
	}
	return data, nil
}

func convertTableProtoToState(id uuid.UUID, data *widgetv1.Table) *state.TableState {
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
		rows := make([]int, len(data.Value.Selection.Rows))
		for i, r := range data.Value.Selection.Rows {
			rows[i] = int(r)
		}
		tableState.Value.Selection = &state.TableStateValueSelection{
			Row:  int(data.Value.Selection.Row),
			Rows: rows,
		}
	}
	return tableState
}
