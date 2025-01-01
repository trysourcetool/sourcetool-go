package sourcetool

import (
	"context"
	"reflect"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/table"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	externaltable "github.com/trysourcetool/sourcetool-go/table"
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
	selection := &table.Selection{
		Row:  0,
		Rows: []int{0},
	}

	state := &table.State{
		ID:           id,
		Data:         data,
		Header:       "Test Table",
		Description:  "Test Description",
		OnSelect:     string(externaltable.OnSelectRerun),
		RowSelection: string(externaltable.RowSelectionSingle),
		Value: table.Value{
			Selection: selection,
		},
	}

	tableData := convertStateToTableData(state)

	if tableData == nil {
		t.Fatal("convertStateToTableData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Header", tableData.Header, state.Header},
		{"Description", tableData.Description, state.Description},
		{"OnSelect", tableData.OnSelect, state.OnSelect},
		{"RowSelection", tableData.RowSelection, state.RowSelection},
		{"Selection.Row", tableData.Value.Selection.Row, state.Value.Selection.Row},
		{"Selection.Rows", tableData.Value.Selection.Rows, state.Value.Selection.Rows},
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
		Header:       "Test Table",
		Description:  "Test Description",
		OnSelect:     externaltable.OnSelectRerun.String(),
		RowSelection: externaltable.RowSelectionSingle.String(),
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
	value := builder.Table(data,
		externaltable.Header(header),
		externaltable.Description(description),
		externaltable.OnSelect(externaltable.OnSelectRerun),
		externaltable.RowSelection(externaltable.RowSelectionMultiple),
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
		{"Header", state.Header, header},
		{"Description", state.Description, description},
		{"OnSelect", state.OnSelect, externaltable.OnSelectRerun.String()},
		{"RowSelection", state.RowSelection, externaltable.RowSelectionMultiple.String()},
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

	// Verify return value
	if !reflect.DeepEqual(value, state.Value) {
		t.Errorf("Return value = %v, want %v", value, state.Value)
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
	if state.OnSelect != externaltable.OnSelectIgnore.String() {
		t.Errorf("Default OnSelect = %v, want %v", state.OnSelect, externaltable.OnSelectIgnore)
	}
	if state.RowSelection != externaltable.RowSelectionSingle.String() {
		t.Errorf("Default RowSelection = %v, want %v", state.RowSelection, externaltable.RowSelectionSingle)
	}
	if state.Header != "" {
		t.Errorf("Default Header = %v, want empty string", state.Header)
	}
	if state.Description != "" {
		t.Errorf("Default Description = %v, want empty string", state.Description)
	}
}
