package sourcetool

import (
	"context"
	"reflect"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/multiselect"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	externalmultiselect "github.com/trysourcetool/sourcetool-go/multiselect"
)

func TestConvertStateToMultiSelectData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := []int{0, 2}
	defaultValue := []int{0}
	options := []string{"Option 1", "Option 2", "Option 3"}

	state := &multiselect.State{
		ID:           id,
		Label:        "Test MultiSelect",
		Value:        value,
		Options:      options,
		Placeholder:  "Select options",
		DefaultValue: defaultValue,
		Required:     true,
		Disabled:     false,
	}

	data := convertStateToMultiSelectData(state)

	if data == nil {
		t.Fatal("convertStateToMultiSelectData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, state.Label},
		{"Value", data.Value, state.Value},
		{"Options length", len(data.Options), len(state.Options)},
		{"Placeholder", data.Placeholder, state.Placeholder},
		{"DefaultValue", data.DefaultValue, state.DefaultValue},
		{"Required", data.Required, state.Required},
		{"Disabled", data.Disabled, state.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertMultiSelectDataToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := []int{0, 2}
	defaultValue := []int{0}
	options := []string{"Option 1", "Option 2", "Option 3"}

	data := &websocket.MultiSelectData{
		Label:        "Test MultiSelect",
		Value:        value,
		Options:      options,
		Placeholder:  "Select options",
		DefaultValue: defaultValue,
		Required:     true,
		Disabled:     false,
	}

	state := convertMultiSelectDataToState(id, data)

	if state == nil {
		t.Fatal("convertMultiSelectDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", state.Value, data.Value},
		{"Options length", len(state.Options), len(data.Options)},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", state.DefaultValue, data.DefaultValue},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestMultiSelect(t *testing.T) {
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

	label := "Test MultiSelect"
	options := []string{"Option 1", "Option 2", "Option 3"}
	defaultValue := []string{"Option 1", "Option 3"}
	placeholder := "Select options"

	// Create MultiSelect component with all options
	value := builder.MultiSelect(label,
		externalmultiselect.Options(options...),
		externalmultiselect.DefaultValue(defaultValue...),
		externalmultiselect.Placeholder(placeholder),
		externalmultiselect.Required(true),
		externalmultiselect.Disabled(true),
	)

	// Verify return value
	if value == nil {
		t.Fatal("MultiSelect returned nil")
	}
	if !reflect.DeepEqual(value.Values, defaultValue) {
		t.Errorf("MultiSelect values = %v, want %v", value.Values, defaultValue)
	}
	expectedIndexes := []int{0, 2}
	if !reflect.DeepEqual(value.Indexes, expectedIndexes) {
		t.Errorf("MultiSelect indexes = %v, want %v", value.Indexes, expectedIndexes)
	}

	// Verify WebSocket message
	if len(mockWS.Messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(mockWS.Messages))
	}
	msg := mockWS.Messages[0]
	if msg.Method != websocket.MessageMethodRenderWidget {
		t.Errorf("WebSocket message method = %v, want %v", msg.Method, websocket.MessageMethodRenderWidget)
	}

	// Verify state
	widgetID := builder.generateMultiSelectID(label, []int{0})
	state := sess.State.GetMultiSelect(widgetID)
	if state == nil {
		t.Fatal("MultiSelect state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Options length", len(state.Options), len(options)},
		{"Placeholder", state.Placeholder, placeholder},
		{"Required", state.Required, true},
		{"Disabled", state.Disabled, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestMultiSelect_WithFormatFunc(t *testing.T) {
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

	label := "Test MultiSelect"
	options := []string{"Option 1", "Option 2"}
	formatFunc := func(value string, index int) string {
		return value + " (Custom)"
	}

	builder.MultiSelect(label,
		externalmultiselect.Options(options...),
		externalmultiselect.FormatFunc(formatFunc),
	)

	// Verify that format function is applied
	widgetID := builder.generateMultiSelectID(label, []int{0})
	state := sess.State.GetMultiSelect(widgetID)
	if state == nil {
		t.Fatal("MultiSelect state not found")
	}

	expectedOptions := []string{"Option 1 (Custom)", "Option 2 (Custom)"}
	if !reflect.DeepEqual(state.Options, expectedOptions) {
		t.Errorf("Formatted options = %v, want %v", state.Options, expectedOptions)
	}
}
