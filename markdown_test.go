package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/markdown"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToMarkdownData(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	state := &markdown.State{
		ID:   id,
		Body: "# Test Markdown",
	}

	data := convertStateToMarkdownData(state)

	if data == nil {
		t.Fatal("convertStateToMarkdownData returned nil")
	}

	if data.Body != state.Body {
		t.Errorf("Body = %v, want %v", data.Body, state.Body)
	}
}

func TestConvertMarkdownDataToState(t *testing.T) {
	data := &websocket.MarkdownData{
		Body: "# Test Markdown",
	}

	state := convertMarkdownDataToState(data)

	if state == nil {
		t.Fatal("convertMarkdownDataToState returned nil")
	}

	if state.Body != data.Body {
		t.Errorf("Body = %v, want %v", state.Body, data.Body)
	}
}

func TestMarkdown(t *testing.T) {
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

	body := "# Test Markdown"
	builder.Markdown(body)

	// Verify WebSocket message
	if len(mockWS.Messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(mockWS.Messages))
	}
	msg := mockWS.Messages[0]
	if msg.Method != websocket.MessageMethodRenderWidget {
		t.Errorf("WebSocket message method = %v, want %v", msg.Method, websocket.MessageMethodRenderWidget)
	}

	// Verify markdown state
	widgetID := builder.generateMarkdownID(body, []int{0})
	state := sess.State.GetMarkdown(widgetID)
	if state == nil {
		t.Fatal("Markdown state not found")
	}

	if state.Body != body {
		t.Errorf("Body = %v, want %v", state.Body, body)
	}
}
