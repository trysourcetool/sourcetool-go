package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	externalcolumns "github.com/trysourcetool/sourcetool-go/columns"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToColumnsData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	columnsState := &state.ColumnsState{
		ID:      id,
		Columns: 3,
	}

	data := convertStateToColumnsData(columnsState)

	if data == nil {
		t.Fatal("convertStateToColumnsData returned nil")
	}

	if data.Columns != columnsState.Columns {
		t.Errorf("Columns = %v, want %v", data.Columns, columnsState.Columns)
	}
}

func TestConvertColumnsDataToState(t *testing.T) {
	data := &websocket.ColumnsData{
		Columns: 3,
	}

	state := convertColumnsDataToState(data)

	if state == nil {
		t.Fatal("convertColumnsDataToState returned nil")
	}

	if state.Columns != data.Columns {
		t.Errorf("Columns = %v, want %v", state.Columns, data.Columns)
	}
}

func TestConvertStateToColumnItemData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	columnItemState := &state.ColumnItemState{
		ID:     id,
		Weight: 0.5,
	}

	data := convertStateToColumnItemData(columnItemState)

	if data == nil {
		t.Fatal("convertStateToColumnItemData returned nil")
	}

	if data.Weight != columnItemState.Weight {
		t.Errorf("Weight = %v, want %v", data.Weight, columnItemState.Weight)
	}
}

func TestConvertColumnItemDataToState(t *testing.T) {
	data := &websocket.ColumnItemData{
		Weight: 0.5,
	}

	state := convertColumnItemDataToState(data)

	if state == nil {
		t.Fatal("convertColumnItemDataToState returned nil")
	}

	if state.Weight != data.Weight {
		t.Errorf("Weight = %v, want %v", state.Weight, data.Weight)
	}
}

func TestColumns(t *testing.T) {
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

	// Test with default options
	cols := 3
	builders := builder.Columns(cols)

	// Verify return value
	if builders == nil {
		t.Fatal("Columns returned nil")
	}
	if len(builders) != cols {
		t.Errorf("Builders length = %v, want %v", len(builders), cols)
	}

	// Verify WebSocket messages
	expectedMessages := cols + 1 // columns widget + column items
	if len(mockWS.Messages) != expectedMessages {
		t.Errorf("WebSocket messages count = %d, want %d", len(mockWS.Messages), expectedMessages)
	}

	// Verify columns state
	widgetID := builder.generateColumnsID([]int{0})
	columnsState := sess.State.GetColumns(widgetID)
	if columnsState == nil {
		t.Fatal("Columns state not found")
	}

	if columnsState.Columns != cols {
		t.Errorf("Columns = %v, want %v", columnsState.Columns, cols)
	}

	// Verify column items state
	for i := 0; i < cols; i++ {
		columnPath := []int{0, i}
		columnID := builder.generateColumnItemID(columnPath)
		columnState := sess.State.Get(columnID)
		if columnState == nil {
			t.Fatalf("Column item state not found for index %d", i)
		}

		expectedWeight := 1.0 / float64(cols)
		columnItemState, ok := columnState.(*state.ColumnItemState)
		if !ok {
			t.Fatalf("Column item state[%d] is not *columnitem.State", i)
		}
		if columnItemState.Weight != expectedWeight {
			t.Errorf("Column item weight[%d] = %v, want %v", i, columnItemState.Weight, expectedWeight)
		}
	}
}

func TestColumns_WithWeight(t *testing.T) {
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

	cols := 3
	weights := []int{2, 1, 1}
	totalWeight := 4

	builders := builder.Columns(cols, externalcolumns.Weight(weights...))

	if builders == nil {
		t.Fatal("Columns returned nil")
	}

	// Verify column items weights
	for i := 0; i < cols; i++ {
		columnPath := []int{0, i}
		columnID := builder.generateColumnItemID(columnPath)
		columnState := sess.State.Get(columnID)
		if columnState == nil {
			t.Fatalf("Column item state not found for index %d", i)
		}

		expectedWeight := float64(weights[i]) / float64(totalWeight)
		columnItemState, ok := columnState.(*state.ColumnItemState)
		if !ok {
			t.Fatalf("Column item state[%d] is not *columnitem.State", i)
		}
		if columnItemState.Weight != expectedWeight {
			t.Errorf("Column item weight[%d] = %v, want %v", i, columnItemState.Weight, expectedWeight)
		}
	}
}

func TestColumns_InvalidInput(t *testing.T) {
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

	tests := []struct {
		name    string
		cols    int
		weights []int
	}{
		{"Zero columns", 0, nil},
		{"Negative columns", -1, nil},
		{"Invalid weights length", 3, []int{1, 1}},
		{"Zero weights", 3, []int{0, 0, 0}},
		{"Negative weights", 3, []int{-1, 1, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builders []UIBuilder
			if tt.weights == nil {
				builders = builder.Columns(tt.cols)
			} else {
				builders = builder.Columns(tt.cols, externalcolumns.Weight(tt.weights...))
			}

			if tt.cols <= 0 && builders != nil {
				t.Error("Expected nil builders for invalid column count")
			}

			if tt.cols > 0 && builders != nil {
				// Verify weights are normalized
				for i := 0; i < tt.cols; i++ {
					columnPath := []int{0, i}
					columnID := builder.generateColumnItemID(columnPath)
					columnState := sess.State.Get(columnID)
					if columnState == nil {
						t.Fatalf("Column item state not found for index %d", i)
					}

					expectedWeight := 1.0 / float64(tt.cols)
					columnItemState, ok := columnState.(*state.ColumnItemState)
					if !ok {
						t.Fatalf("Column item state[%d] is not *columnitem.State", i)
					}
					if columnItemState.Weight != expectedWeight {
						t.Errorf("Column item weight[%d] = %v, want %v", i, columnItemState.Weight, expectedWeight)
					}
				}
			}
		})
	}
}
