package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	"github.com/trysourcetool/sourcetool-go/textinput"
)

func TestConvertStateToTextInputData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	maxLength := 100
	minLength := 10

	textInputState := &state.TextInputState{
		ID:           id,
		Label:        "Test TextInput",
		Value:        "test value",
		Placeholder:  "Enter text",
		DefaultValue: "default",
		Required:     true,
		Disabled:     false,
		MaxLength:    &maxLength,
		MinLength:    &minLength,
	}

	data := convertStateToTextInputData(textInputState)

	if data == nil {
		t.Fatal("convertStateToTextInputData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, textInputState.Label},
		{"Value", data.Value, textInputState.Value},
		{"Placeholder", data.Placeholder, textInputState.Placeholder},
		{"DefaultValue", data.DefaultValue, textInputState.DefaultValue},
		{"Required", data.Required, textInputState.Required},
		{"Disabled", data.Disabled, textInputState.Disabled},
		{"MaxLength", *data.MaxLength, *textInputState.MaxLength},
		{"MinLength", *data.MinLength, *textInputState.MinLength},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertTextInputDataToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	maxLength := 100
	minLength := 10

	data := &websocket.TextInputData{
		Label:        "Test TextInput",
		Value:        "test value",
		Placeholder:  "Enter text",
		DefaultValue: "default",
		Required:     true,
		Disabled:     false,
		MaxLength:    &maxLength,
		MinLength:    &minLength,
	}

	state := convertTextInputDataToState(id, data)

	if state == nil {
		t.Fatal("convertTextInputDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", state.Value, data.Value},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", state.DefaultValue, data.DefaultValue},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
		{"MaxLength", *state.MaxLength, *data.MaxLength},
		{"MinLength", *state.MinLength, *data.MinLength},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestTextInput(t *testing.T) {
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

	label := "Test TextInput"
	defaultValue := "default value"
	placeholder := "Enter text"
	maxLength := 100
	minLength := 10

	// Create TextInput component with all options
	value := builder.TextInput(label,
		textinput.DefaultValue(defaultValue),
		textinput.Placeholder(placeholder),
		textinput.Required(true),
		textinput.Disabled(true),
		textinput.MaxLength(maxLength),
		textinput.MinLength(minLength),
	)

	// Verify return value
	if value != defaultValue {
		t.Errorf("TextInput value = %v, want %v", value, defaultValue)
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
	widgetID := builder.generateTextInputID(label, []int{0})
	state := sess.State.GetTextInput(widgetID)
	if state == nil {
		t.Fatal("TextInput state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Value", state.Value, defaultValue},
		{"Placeholder", state.Placeholder, placeholder},
		{"DefaultValue", state.DefaultValue, defaultValue},
		{"Required", state.Required, true},
		{"Disabled", state.Disabled, true},
		{"MaxLength", *state.MaxLength, maxLength},
		{"MinLength", *state.MinLength, minLength},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
