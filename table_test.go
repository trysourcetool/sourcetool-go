package sourcetool

import (
	"context"
	"reflect"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	"github.com/trysourcetool/sourcetool-go/table"
)

type testData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestConvertStateToTableData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	data := []testData{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}
	selection := &state.TableStateValueSelection{
		Row:  0,
		Rows: []int{0},
	}

	tableState := &state.TableState{
		ID:           id,
		Data:         data,
		Header:       conv.NilValue("Test Table"),
		Description:  conv.NilValue("Test Description"),
		OnSelect:     conv.NilValue(table.SelectionBehaviorRerun.String()),
		RowSelection: conv.NilValue(table.SelectionModeSingle.String()),
		Value: state.TableStateValue{
			Selection: selection,
		},
	}

	tableData := convertStateToTableData(tableState)

	if tableData == nil {
		t.Fatal("convertStateToTableData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Header", conv.SafeValue(tableData.Header), conv.SafeValue(tableState.Header)},
		{"Description", conv.SafeValue(tableData.Description), conv.SafeValue(tableState.Description)},
		{"OnSelect", conv.SafeValue(tableData.OnSelect), conv.SafeValue(tableState.OnSelect)},
		{"RowSelection", tableData.RowSelection, tableState.RowSelection},
		{"Selection.Row", tableData.Value.Selection.Row, tableState.Value.Selection.Row},
		{"Selection.Rows", tableData.Value.Selection.Rows, tableState.Value.Selection.Rows},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	// Verify data separately since it's an any
	if !reflect.DeepEqual(tableData.Data, data) {
		t.Errorf("Data = %v, want %v", tableData.Data, data)
	}
}

func TestConvertTableDataToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	data := []testData{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}
	selection := &websocket.TableDataValueSelection{
		Row:  0,
		Rows: []int{0},
	}

	tableData := &websocket.TableData{
		Data:         data,
		Header:       conv.NilValue("Test Table"),
		Description:  conv.NilValue("Test Description"),
		OnSelect:     conv.NilValue(table.SelectionBehaviorRerun.String()),
		RowSelection: conv.NilValue(table.SelectionModeSingle.String()),
		Value: websocket.TableDataValue{
			Selection: selection,
		},
	}

	state := convertTableDataToState(id, tableData)

	if state == nil {
		t.Fatal("convertTableDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Header", state.Header, tableData.Header},
		{"Description", state.Description, tableData.Description},
		{"OnSelect", state.OnSelect, tableData.OnSelect},
		{"RowSelection", state.RowSelection, tableData.RowSelection},
		{"Selection.Row", state.Value.Selection.Row, tableData.Value.Selection.Row},
		{"Selection.Rows", state.Value.Selection.Rows, tableData.Value.Selection.Rows},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	// Verify data separately since it's an any
	if !reflect.DeepEqual(state.Data, data) {
		t.Errorf("Data = %v, want %v", state.Data, data)
	}
}

func TestTable(t *testing.T) {
	sessionID := uuid.Must(uuid.NewV4())
	pageID := uuid.Must(uuid.NewV4())
	sess := session.New(sessionID, pageID)

	mockWS := mock.NewMockWebSocketClient()

	builder := &uiBuilder{
		context: context.Background(),
		session: sess,
		cursor:  newCursor(),
		page: &page{
			id: pageID,
		},
		runtime: &runtime{
			wsClient: mockWS,
		},
	}

	data := []testData{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}
	header := "Test Table"
	description := "Test Description"

	// Create Table component with all options
	builder.Table(data,
		table.Header(header),
		table.Description(description),
		table.OnSelect(table.SelectionBehaviorRerun),
		table.RowSelection(table.SelectionModeSingle),
	)

	// Verify WebSocket message
	if len(mockWS.Messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(mockWS.Messages))
	}
	msg := mockWS.Messages[0]
	if msg.Method != websocket.MessageMethodRenderWidget {
		t.Errorf("WebSocket message method = %v, want %v", msg.Method, websocket.MessageMethodRenderWidget)
	}

	// Verify state
	widgetID := builder.generateTableID([]int{0})
	state := sess.State.GetTable(widgetID)
	if state == nil {
		t.Fatal("Table state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Header", conv.SafeValue(state.Header), header},
		{"Description", conv.SafeValue(state.Description), description},
		{"OnSelect", conv.SafeValue(state.OnSelect), table.SelectionBehaviorRerun.String()},
		{"RowSelection", conv.SafeValue(state.RowSelection), table.SelectionModeSingle.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	// Verify data separately since it's an any
	if !reflect.DeepEqual(state.Data, data) {
		t.Errorf("Data = %v, want %v", state.Data, data)
	}
}

func TestTable_DefaultValues(t *testing.T) {
	sessionID := uuid.Must(uuid.NewV4())
	pageID := uuid.Must(uuid.NewV4())
	sess := session.New(sessionID, pageID)

	mockWS := mock.NewMockWebSocketClient()

	builder := &uiBuilder{
		context: context.Background(),
		session: sess,
		cursor:  newCursor(),
		page: &page{
			id: pageID,
		},
		runtime: &runtime{
			wsClient: mockWS,
		},
	}

	data := []testData{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}

	// Create Table component without options
	builder.Table(data)

	// Verify state
	widgetID := builder.generateTableID([]int{0})
	state := sess.State.GetTable(widgetID)
	if state == nil {
		t.Fatal("Table state not found")
	}

	// Verify default values
	if conv.SafeValue(state.OnSelect) != table.SelectionBehaviorIgnore.String() {
		t.Errorf("Default OnSelect = %v, want %v", conv.SafeValue(state.OnSelect), table.SelectionBehaviorIgnore)
	}
	if conv.SafeValue(state.RowSelection) != table.SelectionModeSingle.String() {
		t.Errorf("Default RowSelection = %v, want %v", conv.SafeValue(state.RowSelection), table.SelectionModeSingle)
	}
	if state.Header != nil {
		t.Errorf("Default Header = %v, want empty string", state.Header)
	}
	if state.Description != nil {
		t.Errorf("Default Description = %v, want empty string", state.Description)
	}
}
