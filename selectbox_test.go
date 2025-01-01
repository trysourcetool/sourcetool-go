package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/selectbox"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	externalselectbox "github.com/trysourcetool/sourcetool-go/selectbox"
)

func TestConvertStateToSelectboxData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := 1
	defaultValue := 0
	options := []string{"Option 1", "Option 2"}
	placeholder := "Select an option"

	state := &selectbox.State{
		ID:           id,
		Label:        "Test Selectbox",
		Value:        &value,
		Options:      options,
		Placeholder:  placeholder,
		DefaultValue: &defaultValue,
		Required:     true,
		Disabled:     false,
	}

	data := convertStateToSelectboxData(state)

	if data == nil {
		t.Fatal("convertStateToSelectboxData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, state.Label},
		{"Value", *data.Value, *state.Value},
		{"Options length", len(data.Options), len(state.Options)},
		{"Placeholder", data.Placeholder, state.Placeholder},
		{"DefaultValue", *data.DefaultValue, *state.DefaultValue},
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

func TestConvertSelectboxDataToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	value := 1
	defaultValue := 0
	options := []string{"Option 1", "Option 2"}
	placeholder := "Select an option"

	data := &websocket.SelectboxData{
		Label:        "Test Selectbox",
		Value:        &value,
		Options:      options,
		Placeholder:  placeholder,
		DefaultValue: &defaultValue,
		Required:     true,
		Disabled:     false,
	}

	state := convertSelectboxDataToState(id, data)

	if state == nil {
		t.Fatal("convertSelectboxDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Label", state.Label, data.Label},
		{"Value", *state.Value, *data.Value},
		{"Options length", len(state.Options), len(data.Options)},
		{"Placeholder", state.Placeholder, data.Placeholder},
		{"DefaultValue", *state.DefaultValue, *data.DefaultValue},
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

func TestSelectbox(t *testing.T) {
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

	label := "Test Selectbox"
	options := []string{"Option 1", "Option 2"}
	defaultValue := "Option 1"
	placeholder := "Select an option"

	// Create Selectbox component
	value := builder.Selectbox(label,
		externalselectbox.Options(options...),
		externalselectbox.DefaultValue(defaultValue),
		externalselectbox.Placeholder(placeholder),
		externalselectbox.Required(true),
	)

	// Verify return value
	if value == nil {
		t.Fatal("Selectbox returned nil")
	}
	if value.Value != defaultValue {
		t.Errorf("Selectbox value = %v, want %v", value.Value, defaultValue)
	}
	if value.Index != 0 {
		t.Errorf("Selectbox index = %v, want 0", value.Index)
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
	widgetID := builder.generateSelectboxID(label, []int{0})
	state := sess.State.GetSelectbox(widgetID)
	if state == nil {
		t.Fatal("Selectbox state not found")
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
		{"Disabled", state.Disabled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestSelectbox_WithFormatFunc(t *testing.T) {
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

	label := "Test Selectbox"
	options := []string{"Option 1", "Option 2"}
	formatFunc := func(value string, index int) string {
		return value + " (Custom)"
	}

	builder.Selectbox(label,
		externalselectbox.Options(options...),
		externalselectbox.FormatFunc(formatFunc),
	)

	// Verify that format function is applied
	widgetID := builder.generateSelectboxID(label, []int{0})
	state := sess.State.GetSelectbox(widgetID)
	if state == nil {
		t.Fatal("Selectbox state not found")
	}

	expectedOptions := []string{"Option 1 (Custom)", "Option 2 (Custom)"}
	for i, opt := range state.Options {
		if opt != expectedOptions[i] {
			t.Errorf("Formatted option[%d] = %v, want %v", i, opt, expectedOptions[i])
		}
	}
}
