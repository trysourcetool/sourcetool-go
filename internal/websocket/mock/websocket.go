package mock

import "github.com/trysourcetool/sourcetool-go/internal/websocket"

func NewMockWebSocketClient() *mockWebSocketClient {
	return &mockWebSocketClient{}
}

type mockWebSocketClient struct {
	Messages []*websocket.Message
}

func (m *mockWebSocketClient) Enqueue(id string, method websocket.MessageMethod, data any) {
	m.Messages = append(m.Messages, &websocket.Message{
		ID:     id,
		Method: method,
		Kind:   websocket.MessageKindCall,
	})
}

func (m *mockWebSocketClient) EnqueueWithResponse(id string, method websocket.MessageMethod, data any) (*websocket.Message, error) {
	return nil, nil
}

func (m *mockWebSocketClient) RegisterHandler(method websocket.MessageMethod, handler websocket.MessageHandlerFunc) {
}

func (m *mockWebSocketClient) Close() error { return nil }

func (m *mockWebSocketClient) Wait() error { return nil }
