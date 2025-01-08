package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/form"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToFormData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	formState := &state.FormState{
		ID:             id,
		Value:          true,
		ButtonLabel:    "Submit",
		ButtonDisabled: true,
		ClearOnSubmit:  true,
	}

	data := convertStateToFormData(formState)

	if data == nil {
		t.Fatal("convertStateToFormData returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Value", data.Value, formState.Value},
		{"ButtonLabel", data.ButtonLabel, formState.ButtonLabel},
		{"ButtonDisabled", data.ButtonDisabled, formState.ButtonDisabled},
		{"ClearOnSubmit", data.ClearOnSubmit, formState.ClearOnSubmit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConvertFormDataToState(t *testing.T) {
	data := &websocket.FormData{
		Value:          true,
		ButtonLabel:    "Submit",
		ButtonDisabled: true,
		ClearOnSubmit:  true,
	}

	state := convertFormDataToState(data)

	if state == nil {
		t.Fatal("convertFormDataToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Value", state.Value, data.Value},
		{"ButtonLabel", state.ButtonLabel, data.ButtonLabel},
		{"ButtonDisabled", state.ButtonDisabled, data.ButtonDisabled},
		{"ClearOnSubmit", state.ClearOnSubmit, data.ClearOnSubmit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestForm(t *testing.T) {
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

	buttonLabel := "Submit"
	childBuilder, submitted := builder.Form(buttonLabel)

	// Verify return values
	if childBuilder == nil {
		t.Fatal("Form returned nil builder")
	}
	if submitted {
		t.Error("Form returned true for submitted, want false")
	}

	// Verify WebSocket message
	if len(mockWS.Messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(mockWS.Messages))
	}
	msg := mockWS.Messages[0]
	if msg.Method != websocket.MessageMethodRenderWidget {
		t.Errorf("WebSocket message method = %v, want %v", msg.Method, websocket.MessageMethodRenderWidget)
	}

	// Verify form state
	widgetID := builder.generateFormID([]int{0})
	state := sess.State.GetForm(widgetID)
	if state == nil {
		t.Fatal("Form state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ButtonLabel", state.ButtonLabel, buttonLabel},
		{"ButtonDisabled", state.ButtonDisabled, false}, // default value
		{"ClearOnSubmit", state.ClearOnSubmit, false},   // default value
		{"Value", state.Value, false},                   // default value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestForm_WithOptions(t *testing.T) {
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

	buttonLabel := "Submit"
	childBuilder, submitted := builder.Form(buttonLabel,
		form.ButtonDisabled(true),
		form.ClearOnSubmit(true),
	)

	if childBuilder == nil {
		t.Fatal("Form returned nil builder")
	}
	if submitted {
		t.Error("Form returned true for submitted, want false")
	}

	// Verify form state with options
	widgetID := builder.generateFormID([]int{0})
	state := sess.State.GetForm(widgetID)
	if state == nil {
		t.Fatal("Form state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ButtonLabel", state.ButtonLabel, buttonLabel},
		{"ButtonDisabled", state.ButtonDisabled, true},
		{"ClearOnSubmit", state.ClearOnSubmit, true},
		{"Value", state.Value, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
