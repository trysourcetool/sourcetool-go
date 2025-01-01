package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/textarea"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	externaltextarea "github.com/trysourcetool/sourcetool-go/textarea"
)

func TestConvertStateToTextAreaData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	maxLength := 1000
	minLength := 10
	maxLines := 10
	minLines := 3

	state := &textarea.State{
		ID:           id,
		Label:        "Test TextArea",
		Value:        "test value",
		Placeholder:  "Enter text",
		DefaultValue: "default",
		Required:     true,
		Disabled:     false,
		MaxLength:    &maxLength,
		MinLength:    &minLength,
		MaxLines:     &maxLines,
		MinLines:     &minLines,
		AutoResize:   true,
	}

	data := convertStateToTextAreaData(state)

	if data == nil {
		t.Fatal("convertStateToTextAreaData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Label", data.Label, state.Label},
		{"Value", data.Value, state.Value},
		{"Placeholder", data.Placeholder, state.Placeholder},
		{"DefaultValue", data.DefaultValue, state.DefaultValue},
		{"Required", data.Required, state.Required},
		{"Disabled", data.Disabled, state.Disabled},
		{"MaxLength", *data.MaxLength, *state.MaxLength},
		{"MinLength", *data.MinLength, *state.MinLength},
		{"MaxLines", *data.MaxLines, *state.MaxLines},
		{"MinLines", *data.MinLines, *state.MinLines},
		{"AutoResize", data.AutoResize, state.AutoResize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertTextAreaDataToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	maxLength := 1000
	minLength := 10
	maxLines := 10
	minLines := 3

	data := &websocket.TextAreaData{
		Label:        "Test TextArea",
		Value:        "test value",
		Placeholder:  "Enter text",
		DefaultValue: "default",
		Required:     true,
		Disabled:     false,
		MaxLength:    &maxLength,
		MinLength:    &minLength,
		MaxLines:     &maxLines,
		MinLines:     &minLines,
		AutoResize:   true,
	}

	state := convertTextAreaDataToState(id, data)

	if state == nil {
		t.Fatal("convertTextAreaDataToState returned nil")
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
		{"MaxLines", *state.MaxLines, *data.MaxLines},
		{"MinLines", *state.MinLines, *data.MinLines},
		{"AutoResize", state.AutoResize, data.AutoResize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestTextArea(t *testing.T) {
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

	label := "Test TextArea"
	defaultValue := "default value"
	placeholder := "Enter text"
	maxLength := 1000
	minLength := 10
	maxLines := 10
	minLines := 3

	// Create TextArea component with all options
	value := builder.TextArea(label,
		externaltextarea.DefaultValue(defaultValue),
		externaltextarea.Placeholder(placeholder),
		externaltextarea.Required(true),
		externaltextarea.Disabled(true),
		externaltextarea.MaxLength(maxLength),
		externaltextarea.MinLength(minLength),
		externaltextarea.MaxLines(maxLines),
		externaltextarea.MinLines(minLines),
		externaltextarea.AutoResize(false),
	)

	// Verify return value
	if value != defaultValue {
		t.Errorf("TextArea value = %v, want %v", value, defaultValue)
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
	widgetID := builder.generateTextAreaID(label, []int{0})
	state := sess.State.GetTextArea(widgetID)
	if state == nil {
		t.Fatal("TextArea state not found")
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
		{"MaxLines", *state.MaxLines, maxLines},
		{"MinLines", *state.MinLines, minLines},
		{"AutoResize", state.AutoResize, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestTextArea_DefaultMinLines(t *testing.T) {
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

	label := "Test TextArea"

	// Create TextArea component without options
	builder.TextArea(label)

	// Verify state
	widgetID := builder.generateTextAreaID(label, []int{0})
	state := sess.State.GetTextArea(widgetID)
	if state == nil {
		t.Fatal("TextArea state not found")
	}

	// Verify default MinLines value
	if state.MinLines == nil {
		t.Fatal("MinLines is nil, want 2")
	}
	if *state.MinLines != 2 {
		t.Errorf("Default MinLines = %v, want 2", *state.MinLines)
	}
}
