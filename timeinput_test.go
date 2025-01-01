package sourcetool

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/timeinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	externaltimeinput "github.com/trysourcetool/sourcetool-go/timeinput"
)

func TestConvertStateToTimeInputData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	now := time.Now()

	state := &timeinput.State{
		ID:           id,
		Label:        "Test TimeInput",
		Value:        &now,
		Placeholder:  "Select time",
		DefaultValue: &now,
		Required:     true,
		Disabled:     false,
		Location:     time.Local,
	}

	data := convertStateToTimeInputData(state)

	if data == nil {
		t.Fatal("convertStateToTimeInputData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, state.Label},
		{"Value", data.Value, state.Value.Format(time.TimeOnly)},
		{"Placeholder", data.Placeholder, state.Placeholder},
		{"DefaultValue", data.DefaultValue, state.DefaultValue.Format(time.TimeOnly)},
		{"Required", data.Required, state.Required},
		{"Disabled", data.Disabled, state.Disabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertTimeInputDataToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	now := time.Now()
	timeStr := now.Format(time.TimeOnly)

	data := &websocket.TimeInputData{
		Label:        "Test TimeInput",
		Value:        timeStr,
		Placeholder:  "Select time",
		DefaultValue: timeStr,
		Required:     true,
		Disabled:     false,
	}

	state, err := convertTimeInputDataToState(id, data, time.Local)
	if err != nil {
		t.Fatalf("convertTimeInputDataToState returned error: %v", err)
	}

	if state == nil {
		t.Fatal("convertTimeInputDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", state.Value.Format(time.TimeOnly), data.Value},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", state.DefaultValue.Format(time.TimeOnly), data.DefaultValue},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
		{"Location", state.Location.String(), time.Local.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertTimeInputDataToState_InvalidTime(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	data := &websocket.TimeInputData{
		Value: "invalid-time",
	}

	_, err := convertTimeInputDataToState(id, data, time.Local)
	if err == nil {
		t.Error("Expected error for invalid time, got nil")
	}
}

func TestTimeInput(t *testing.T) {
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

	label := "Test TimeInput"
	now := time.Now()
	placeholder := "Select time"
	location := *time.UTC

	// Create TimeInput component with all options
	value := builder.TimeInput(label,
		externaltimeinput.DefaultValue(now),
		externaltimeinput.Placeholder(placeholder),
		externaltimeinput.Required(true),
		externaltimeinput.Disabled(true),
		externaltimeinput.Location(location),
	)

	// Verify return value
	if value == nil {
		t.Fatal("TimeInput returned nil")
	}
	if !value.Equal(now) {
		t.Errorf("TimeInput value = %v, want %v", value, now)
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
	widgetID := builder.generateTimeInputID(label, []int{0})
	state := sess.State.GetTimeInput(widgetID)
	if state == nil {
		t.Fatal("TimeInput state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", state.Label, label},
		{"Placeholder", state.Placeholder, placeholder},
		{"Required", state.Required, true},
		{"Disabled", state.Disabled, true},
		{"Location", state.Location.String(), location.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	// Verify time values separately to handle time comparison
	if !state.Value.Equal(*value) {
		t.Errorf("Value = %v, want %v", state.Value, value)
	}
	if !state.DefaultValue.Equal(now) {
		t.Errorf("DefaultValue = %v, want %v", state.DefaultValue, now)
	}
}

func TestTimeInput_DefaultValues(t *testing.T) {
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

	label := "Test TimeInput"

	// Create TimeInput component without options
	builder.TimeInput(label)

	// Verify state
	widgetID := builder.generateTimeInputID(label, []int{0})
	state := sess.State.GetTimeInput(widgetID)
	if state == nil {
		t.Fatal("TimeInput state not found")
	}

	// Verify default values
	if state.Location != time.Local {
		t.Errorf("Default Location = %v, want %v", state.Location, time.Local)
	}
	if state.Required {
		t.Error("Default Required = true, want false")
	}
	if state.Disabled {
		t.Error("Default Disabled = true, want false")
	}
}
