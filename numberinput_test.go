package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	"github.com/trysourcetool/sourcetool-go/numberinput"
)

func TestConvertStateToNumberInputData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := 42.5
	defaultValue := 0.0
	maxValue := 100.0
	minValue := 0.0

	numberInputState := &state.NumberInputState{
		ID:           id,
		Label:        "Test NumberInput",
		Value:        &value,
		Placeholder:  "Enter number",
		DefaultValue: &defaultValue,
		Required:     true,
		Disabled:     false,
		MaxValue:     &maxValue,
		MinValue:     &minValue,
	}

	data := convertStateToNumberInputData(numberInputState)

	if data == nil {
		t.Fatal("convertStateToNumberInputData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, numberInputState.Label},
		{"Value", *data.Value, *numberInputState.Value},
		{"Placeholder", data.Placeholder, numberInputState.Placeholder},
		{"DefaultValue", *data.DefaultValue, *numberInputState.DefaultValue},
		{"Required", data.Required, numberInputState.Required},
		{"Disabled", data.Disabled, numberInputState.Disabled},
		{"MaxValue", *data.MaxValue, *numberInputState.MaxValue},
		{"MinValue", *data.MinValue, *numberInputState.MinValue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertNumberInputDataToState(t *testing.T) {
	value := 42.5
	defaultValue := 0.0
	maxValue := 100.0
	minValue := 0.0

	data := &websocket.NumberInputData{
		Label:        "Test NumberInput",
		Value:        &value,
		Placeholder:  "Enter number",
		DefaultValue: &defaultValue,
		Required:     true,
		Disabled:     false,
		MaxValue:     &maxValue,
		MinValue:     &minValue,
	}

	state := convertNumberInputDataToState(data)

	if state == nil {
		t.Fatal("convertNumberInputDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, data.Label},
		{"Value", *state.Value, *data.Value},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", *state.DefaultValue, *data.DefaultValue},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
		{"MaxValue", *state.MaxValue, *data.MaxValue},
		{"MinValue", *state.MinValue, *data.MinValue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestNumberInput(t *testing.T) {
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

	label := "Test NumberInput"
	defaultValue := 42.5
	placeholder := "Enter number"
	maxValue := 100.0
	minValue := 0.0

	// Create NumberInput component with all options
	value := builder.NumberInput(label,
		numberinput.DefaultValue(defaultValue),
		numberinput.Placeholder(placeholder),
		numberinput.Required(true),
		numberinput.Disabled(true),
		numberinput.MaxValue(maxValue),
		numberinput.MinValue(minValue),
	)

	// Verify return value
	if value == nil {
		t.Fatal("NumberInput returned nil")
	}
	if *value != defaultValue {
		t.Errorf("NumberInput value = %v, want %v", *value, defaultValue)
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
	widgetID := builder.generateNumberInputID(label, []int{0})
	state := sess.State.GetNumberInput(widgetID)
	if state == nil {
		t.Fatal("NumberInput state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Value", *state.Value, defaultValue},
		{"Placeholder", state.Placeholder, placeholder},
		{"DefaultValue", *state.DefaultValue, defaultValue},
		{"Required", state.Required, true},
		{"Disabled", state.Disabled, true},
		{"MaxValue", *state.MaxValue, maxValue},
		{"MinValue", *state.MinValue, minValue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
