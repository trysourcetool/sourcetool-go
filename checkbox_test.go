package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/checkbox"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToCheckboxData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())

	checkboxState := &state.CheckboxState{
		ID:           id,
		Label:        "Test Checkbox",
		Value:        true,
		DefaultValue: false,
		Required:     true,
		Disabled:     false,
	}

	data := convertStateToCheckboxData(checkboxState)

	if data == nil {
		t.Fatal("convertStateToCheckboxData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, checkboxState.Label},
		{"Value", data.Value, checkboxState.Value},
		{"DefaultValue", data.DefaultValue, checkboxState.DefaultValue},
		{"Required", data.Required, checkboxState.Required},
		{"Disabled", data.Disabled, checkboxState.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertCheckboxDataToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())

	data := &websocket.CheckboxData{
		Label:        "Test Checkbox",
		Value:        true,
		DefaultValue: false,
		Required:     true,
		Disabled:     false,
	}

	state := convertCheckboxDataToState(id, data)

	if state == nil {
		t.Fatal("convertCheckboxDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", state.Value, data.Value},
		{"DefaultValue", state.DefaultValue, data.DefaultValue},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestCheckbox(t *testing.T) {
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

	label := "Test Checkbox"

	// Create Checkbox component with all options
	value := builder.Checkbox(label,
		checkbox.DefaultValue(true),
		checkbox.Required(true),
		checkbox.Disabled(true),
	)

	// Verify return value
	if !value {
		t.Error("Checkbox value = false, want true")
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
	widgetID := builder.generateCheckboxID(label, []int{0})
	state := sess.State.GetCheckbox(widgetID)
	if state == nil {
		t.Fatal("Checkbox state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Value", state.Value, true},
		{"DefaultValue", state.DefaultValue, true},
		{"Required", state.Required, true},
		{"Disabled", state.Disabled, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
