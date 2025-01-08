package sourcetool

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/datetimeinput"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToDateTimeInputData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	now := time.Now()
	maxDate := now.AddDate(1, 0, 0)
	minDate := now.AddDate(-1, 0, 0)

	dateTimeInputState := &state.DateTimeInputState{
		ID:           id,
		Label:        "Test DateTimeInput",
		Value:        &now,
		Placeholder:  "Select date and time",
		DefaultValue: &now,
		Required:     true,
		Disabled:     false,
		Format:       "YYYY/MM/DD HH:MM:SS",
		MaxValue:     &maxDate,
		MinValue:     &minDate,
		Location:     time.Local,
	}

	data := convertStateToDateTimeInputData(dateTimeInputState)

	if data == nil {
		t.Fatal("convertStateToDateTimeInputData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, dateTimeInputState.Label},
		{"Value", data.Value, dateTimeInputState.Value.Format(time.DateTime)},
		{"Placeholder", data.Placeholder, dateTimeInputState.Placeholder},
		{"DefaultValue", data.DefaultValue, dateTimeInputState.DefaultValue.Format(time.DateTime)},
		{"Required", data.Required, dateTimeInputState.Required},
		{"Disabled", data.Disabled, dateTimeInputState.Disabled},
		{"Format", data.Format, dateTimeInputState.Format},
		{"MaxValue", data.MaxValue, dateTimeInputState.MaxValue.Format(time.DateTime)},
		{"MinValue", data.MinValue, dateTimeInputState.MinValue.Format(time.DateTime)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertDateTimeInputDataToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	now := time.Now()
	dateStr := now.Format(time.DateTime)
	maxDateStr := now.AddDate(1, 0, 0).Format(time.DateTime)
	minDateStr := now.AddDate(-1, 0, 0).Format(time.DateTime)

	data := &websocket.DateTimeInputData{
		Label:        "Test DateTimeInput",
		Value:        dateStr,
		Placeholder:  "Select date and time",
		DefaultValue: dateStr,
		Required:     true,
		Disabled:     false,
		Format:       "YYYY/MM/DD HH:MM:SS",
		MaxValue:     maxDateStr,
		MinValue:     minDateStr,
	}

	state, err := convertDateTimeInputDataToState(id, data, time.Local)
	if err != nil {
		t.Fatalf("convertDateTimeInputDataToState returned error: %v", err)
	}

	if state == nil {
		t.Fatal("convertDateTimeInputDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", state.Value.Format(time.DateTime), data.Value},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", state.DefaultValue.Format(time.DateTime), data.DefaultValue},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
		{"Format", state.Format, data.Format},
		{"MaxValue", state.MaxValue.Format(time.DateTime), data.MaxValue},
		{"MinValue", state.MinValue.Format(time.DateTime), data.MinValue},
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

func TestConvertDateTimeInputDataToState_InvalidDateTime(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	data := &websocket.DateTimeInputData{
		Value: "invalid-datetime",
	}

	_, err := convertDateTimeInputDataToState(id, data, time.Local)
	if err == nil {
		t.Error("Expected error for invalid datetime, got nil")
	}
}

func TestDateTimeInput(t *testing.T) {
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

	label := "Test DateTimeInput"
	now := time.Now()
	maxDate := now.AddDate(1, 0, 0)
	minDate := now.AddDate(-1, 0, 0)
	placeholder := "Select date and time"
	format := "YYYY-MM-DD HH:mm:ss"
	location := *time.UTC

	// Create DateTimeInput component with all options
	value := builder.DateTimeInput(label,
		datetimeinput.DefaultValue(now),
		datetimeinput.Placeholder(placeholder),
		datetimeinput.Required(true),
		datetimeinput.Disabled(true),
		datetimeinput.Format(format),
		datetimeinput.MaxValue(maxDate),
		datetimeinput.MinLength(minDate),
		datetimeinput.Location(location),
	)

	// Verify return value
	if value == nil {
		t.Fatal("DateTimeInput returned nil")
	}
	if !value.Equal(now) {
		t.Errorf("DateTimeInput value = %v, want %v", value, now)
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
	widgetID := builder.generateDateTimeInputID(label, []int{0})
	state := sess.State.GetDateTimeInput(widgetID)
	if state == nil {
		t.Fatal("DateTimeInput state not found")
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
		{"Format", state.Format, format},
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
	if !state.MaxValue.Equal(maxDate) {
		t.Errorf("MaxValue = %v, want %v", state.MaxValue, maxDate)
	}
	if !state.MinValue.Equal(minDate) {
		t.Errorf("MinValue = %v, want %v", state.MinValue, minDate)
	}
}

func TestDateTimeInput_DefaultValues(t *testing.T) {
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

	label := "Test DateTimeInput"

	// Create DateTimeInput component without options
	builder.DateTimeInput(label)

	// Verify state
	widgetID := builder.generateDateTimeInputID(label, []int{0})
	state := sess.State.GetDateTimeInput(widgetID)
	if state == nil {
		t.Fatal("DateTimeInput state not found")
	}

	// Verify default values
	if state.Format != "YYYY/MM/DD HH:MM:SS" {
		t.Errorf("Default Format = %v, want YYYY/MM/DD HH:MM:SS", state.Format)
	}
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
