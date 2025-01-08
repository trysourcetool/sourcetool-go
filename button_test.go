package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	externalbutton "github.com/trysourcetool/sourcetool-go/button"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToButtonData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())

	buttonState := &state.ButtonState{
		ID:       id,
		Label:    "Test Button",
		Value:    true,
		Disabled: true,
	}

	data := convertStateToButtonData(buttonState)

	if data == nil {
		t.Fatal("convertStateToButtonData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, buttonState.Label},
		{"Value", data.Value, buttonState.Value},
		{"Disabled", data.Disabled, buttonState.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertButtonDataToState(t *testing.T) {
	data := &websocket.ButtonData{
		Label:    "Test Button",
		Value:    true,
		Disabled: true,
	}

	state := convertButtonDataToState(data)

	if state == nil {
		t.Fatal("convertButtonDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, data.Label},
		{"Value", state.Value, data.Value},
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

func TestButton(t *testing.T) {
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

	label := "Test Button"

	// Create Button component with all options
	value := builder.Button(label,
		externalbutton.Disabled(true),
	)

	// Verify return value
	if value {
		t.Error("Button value = true, want false")
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
	widgetID := builder.generateButtonInputID(label, []int{0})
	state := sess.State.GetButton(widgetID)
	if state == nil {
		t.Fatal("Button state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Value", state.Value, false},
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

func TestButton_DefaultState(t *testing.T) {
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

	label := "Test Button"

	// Create Button component without options
	builder.Button(label)

	// Verify state
	widgetID := builder.generateButtonInputID(label, []int{0})
	state := sess.State.GetButton(widgetID)
	if state == nil {
		t.Fatal("Button state not found")
	}

	// Verify default values
	if state.Value {
		t.Error("Default Value = true, want false")
	}
	if state.Disabled {
		t.Error("Default Disabled = true, want false")
	}
}
