package sourcetool

import (
	"context"
	"reflect"
	"testing"

	"github.com/gofrs/uuid/v5"

	externalcheckboxgroup "github.com/trysourcetool/sourcetool-go/checkboxgroup"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToCheckboxGroupData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := []int{0, 2}
	defaultValue := []int{0}
	options := []string{"Option 1", "Option 2", "Option 3"}

	checkboxGroupState := &state.CheckboxGroupState{
		ID:           id,
		Label:        "Test CheckboxGroup",
		Value:        value,
		Options:      options,
		DefaultValue: defaultValue,
		Required:     true,
		Disabled:     false,
	}

	data := convertStateToCheckboxGroupData(checkboxGroupState)

	if data == nil {
		t.Fatal("convertStateToCheckboxGroupData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, checkboxGroupState.Label},
		{"Value", data.Value, checkboxGroupState.Value},
		{"Options length", len(data.Options), len(checkboxGroupState.Options)},
		{"DefaultValue", data.DefaultValue, checkboxGroupState.DefaultValue},
		{"Required", data.Required, checkboxGroupState.Required},
		{"Disabled", data.Disabled, checkboxGroupState.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertCheckboxGroupDataToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := []int{0, 2}
	defaultValue := []int{0}
	options := []string{"Option 1", "Option 2", "Option 3"}

	data := &websocket.CheckboxGroupData{
		Label:        "Test CheckboxGroup",
		Value:        value,
		Options:      options,
		DefaultValue: defaultValue,
		Required:     true,
		Disabled:     false,
	}

	state := convertCheckboxGroupDataToState(id, data)

	if state == nil {
		t.Fatal("convertCheckboxGroupDataToState returned nil")
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

func TestCheckboxGroup(t *testing.T) {
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

	label := "Test CheckboxGroup"
	options := []string{"Option 1", "Option 2", "Option 3"}
	defaultValue := []string{"Option 1", "Option 3"}

	// Create CheckboxGroup component with all options
	value := builder.CheckboxGroup(label,
		externalcheckboxgroup.Options(options...),
		externalcheckboxgroup.DefaultValue(defaultValue...),
		externalcheckboxgroup.Required(true),
		externalcheckboxgroup.Disabled(true),
	)

	// Verify return value
	if value == nil {
		t.Fatal("CheckboxGroup returned nil")
	}
	if !reflect.DeepEqual(value.Values, defaultValue) {
		t.Errorf("CheckboxGroup values = %v, want %v", value.Values, defaultValue)
	}
	expectedIndexes := []int{0, 2}
	if !reflect.DeepEqual(value.Indexes, expectedIndexes) {
		t.Errorf("CheckboxGroup indexes = %v, want %v", value.Indexes, expectedIndexes)
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
	widgetID := builder.generateCheckboxGroupID(label, []int{0})
	state := sess.State.GetCheckboxGroup(widgetID)
	if state == nil {
		t.Fatal("CheckboxGroup state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Options length", len(state.Options), len(options)},
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

func TestCheckboxGroup_WithFormatFunc(t *testing.T) {
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

	label := "Test CheckboxGroup"
	options := []string{"Option 1", "Option 2"}
	formatFunc := func(value string, index int) string {
		return value + " (Custom)"
	}

	builder.CheckboxGroup(label,
		externalcheckboxgroup.Options(options...),
		externalcheckboxgroup.FormatFunc(formatFunc),
	)

	// Verify that format function is applied
	widgetID := builder.generateCheckboxGroupID(label, []int{0})
	state := sess.State.GetCheckboxGroup(widgetID)
	if state == nil {
		t.Fatal("CheckboxGroup state not found")
	}

	expectedOptions := []string{"Option 1 (Custom)", "Option 2 (Custom)"}
	if !reflect.DeepEqual(state.Options, expectedOptions) {
		t.Errorf("Formatted options = %v, want %v", state.Options, expectedOptions)
	}
}
