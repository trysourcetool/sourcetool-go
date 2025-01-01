package sourcetool

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestRuntime_HandleInitializeClient(t *testing.T) {
	pages := make(map[uuid.UUID]*page)
	pageID := uuid.Must(uuid.NewV4())

	// Test page handler
	handlerCalled := false
	testPage := &page{
		id:   pageID,
		name: "Test Page",
		handler: func(ui UIBuilder) error {
			handlerCalled = true
			return nil
		},
	}
	pages[pageID] = testPage

	r := &runtime{
		wsClient:       mock.NewMockWebSocketClient(),
		sessionManager: session.NewSessionManager(),
		pageManager:    newPageManager(pages),
	}

	// Create test message
	sessionID := uuid.Must(uuid.NewV4())
	payload := websocket.InitializeClientPayload{
		SessionID: sessionID.String(),
		PageID:    pageID.String(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	msg := &websocket.Message{
		ID:      uuid.Must(uuid.NewV4()).String(),
		Kind:    websocket.MessageKindCall,
		Method:  websocket.MessageMethodInitializeClient,
		Payload: payloadBytes,
	}

	// Execute handler
	if err := r.handleInitializeCilent(msg); err != nil {
		t.Fatalf("handleInitializeClient failed: %v", err)
	}

	// Verify that session was created
	sess := r.sessionManager.GetSession(sessionID)
	if sess == nil {
		t.Error("session was not created")
	}

	// Verify that page handler was called
	if !handlerCalled {
		t.Error("page handler was not called")
	}
}

func TestRuntime_HandleRerunPage(t *testing.T) {
	pages := make(map[uuid.UUID]*page)
	pageID := uuid.Must(uuid.NewV4())
	sessionID := uuid.Must(uuid.NewV4())

	// Test page handler
	handlerCallCount := 0
	testPage := &page{
		id:   pageID,
		name: "Test Page",
		handler: func(ui UIBuilder) error {
			handlerCallCount++
			return nil
		},
	}
	pages[pageID] = testPage

	r := &runtime{
		wsClient:       mock.NewMockWebSocketClient(),
		sessionManager: session.NewSessionManager(),
		pageManager:    newPageManager(pages),
	}

	// Initialize session
	sess := session.New(sessionID, pageID)
	r.sessionManager.SetSession(sess)

	// Create test message
	payload := websocket.RerunPagePayload{
		SessionID: sessionID.String(),
		PageID:    pageID.String(),
		State:     json.RawMessage("{}"),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	msg := &websocket.Message{
		ID:      uuid.Must(uuid.NewV4()).String(),
		Kind:    websocket.MessageKindCall,
		Method:  websocket.MessageMethodRerunPage,
		Payload: payloadBytes,
	}

	// Execute handler
	if err := r.handleRerunPage(msg); err != nil {
		t.Fatalf("handleRerunPage failed: %v", err)
	}

	// Verify that page handler was called
	if handlerCallCount != 1 {
		t.Errorf("page handler call count = %d, want 1", handlerCallCount)
	}

	// Execute again
	if err := r.handleRerunPage(msg); err != nil {
		t.Fatalf("second handleRerunPage failed: %v", err)
	}

	// Verify that page handler was called again
	if handlerCallCount != 2 {
		t.Errorf("page handler call count after rerun = %d, want 2", handlerCallCount)
	}
}

func TestRuntime_HandleCloseSession(t *testing.T) {
	sessionID := uuid.Must(uuid.NewV4())
	pageID := uuid.Must(uuid.NewV4())

	r := &runtime{
		wsClient:       mock.NewMockWebSocketClient(),
		sessionManager: session.NewSessionManager(),
		pageManager:    newPageManager(make(map[uuid.UUID]*page)),
	}

	// Initialize session
	sess := session.New(sessionID, pageID)
	r.sessionManager.SetSession(sess)

	// Create test message
	payload := websocket.CloseSessionPayload{
		SessionID: sessionID.String(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	msg := &websocket.Message{
		ID:      uuid.Must(uuid.NewV4()).String(),
		Kind:    websocket.MessageKindCall,
		Method:  websocket.MessageMethodCloseSession,
		Payload: payloadBytes,
	}

	// Execute handler
	if err := r.handleCloseSession(msg); err != nil {
		t.Fatalf("handleCloseSession failed: %v", err)
	}

	// Verify that session was deleted
	if got := r.sessionManager.GetSession(sessionID); got != nil {
		t.Error("session was not deleted")
	}
}
