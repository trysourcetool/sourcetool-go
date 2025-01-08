package sourcetool

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"

	externaldateinput "github.com/trysourcetool/sourcetool-go/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToDateInputData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	now := time.Now()
	maxDate := now.AddDate(1, 0, 0)
	minDate := now.AddDate(-1, 0, 0)

	dateInputState := &state.DateInputState{
		ID:           id,
		Label:        "Test DateInput",
		Value:        &now,
		Placeholder:  "Select date",
		DefaultValue: &now,
		Required:     true,
		Disabled:     false,
		Format:       "YYYY/MM/DD",
		MaxValue:     &maxDate,
		MinValue:     &minDate,
		Location:     time.Local,
	}

	data := convertStateToDateInputData(dateInputState)

	if data == nil {
		t.Fatal("convertStateToDateInputData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, dateInputState.Label},
		{"Value", data.Value, dateInputState.Value.Format(time.DateOnly)},
		{"Placeholder", data.Placeholder, dateInputState.Placeholder},
		{"DefaultValue", data.DefaultValue, dateInputState.DefaultValue.Format(time.DateOnly)},
		{"Required", data.Required, dateInputState.Required},
		{"Disabled", data.Disabled, dateInputState.Disabled},
		{"Format", data.Format, dateInputState.Format},
		{"MaxValue", data.MaxValue, dateInputState.MaxValue.Format(time.DateOnly)},
		{"MinValue", data.MinValue, dateInputState.MinValue.Format(time.DateOnly)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertDateInputDataToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	now := time.Now()
	dateStr := now.Format(time.DateOnly)
	maxDateStr := now.AddDate(1, 0, 0).Format(time.DateOnly)
	minDateStr := now.AddDate(-1, 0, 0).Format(time.DateOnly)

	data := &websocket.DateInputData{
		Label:        "Test DateInput",
		Value:        dateStr,
		Placeholder:  "Select date",
		DefaultValue: dateStr,
		Required:     true,
		Disabled:     false,
		Format:       "YYYY/MM/DD",
		MaxValue:     maxDateStr,
		MinValue:     minDateStr,
	}

	state, err := convertDateInputDataToState(id, data, time.Local)
	if err != nil {
		t.Fatalf("convertDateInputDataToState returned error: %v", err)
	}

	if state == nil {
		t.Fatal("convertDateInputDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", state.Value.Format(time.DateOnly), data.Value},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", state.DefaultValue.Format(time.DateOnly), data.DefaultValue},
		{"Required", state.Required, data.Required},
		{"Disabled", state.Disabled, data.Disabled},
		{"Format", state.Format, data.Format},
		{"MaxValue", state.MaxValue.Format(time.DateOnly), data.MaxValue},
		{"MinValue", state.MinValue.Format(time.DateOnly), data.MinValue},
		{"Location", state.Location, time.Local},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertDateInputDataToState_InvalidDate(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	data := &websocket.DateInputData{
		Value: "invalid-date",
	}

	_, err := convertDateInputDataToState(id, data, time.Local)
	if err == nil {
		t.Error("Expected error for invalid date, got nil")
	}
}

func TestDateInput(t *testing.T) {
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

	label := "Test DateInput"
	now := time.Now()
	maxDate := now.AddDate(1, 0, 0)
	minDate := now.AddDate(-1, 0, 0)
	placeholder := "Select date"
	format := "YYYY-MM-DD"
	location := *time.UTC

	// Create DateInput component with all options
	value := builder.DateInput(label,
		externaldateinput.DefaultValue(now),
		externaldateinput.Placeholder(placeholder),
		externaldateinput.Required(true),
		externaldateinput.Disabled(true),
		externaldateinput.Format(format),
		externaldateinput.MaxValue(maxDate),
		externaldateinput.MinLength(minDate),
		externaldateinput.Location(location),
	)

	// Verify return value
	if value == nil {
		t.Fatal("DateInput returned nil")
	}
	if !value.Equal(now) {
		t.Errorf("DateInput value = %v, want %v", value, now)
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
	widgetID := builder.generateDateInputID(label, []int{0})
	state := sess.State.GetDateInput(widgetID)
	if state == nil {
		t.Fatal("DateInput state not found")
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

func TestDateInput_DefaultValues(t *testing.T) {
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

	label := "Test DateInput"

	// Create DateInput component without options
	builder.DateInput(label)

	// Verify state
	widgetID := builder.generateDateInputID(label, []int{0})
	state := sess.State.GetDateInput(widgetID)
	if state == nil {
		t.Fatal("DateInput state not found")
	}

	// Verify default values
	if state.Format != "YYYY/MM/DD" {
		t.Errorf("Default Format = %v, want YYYY/MM/DD", state.Format)
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
