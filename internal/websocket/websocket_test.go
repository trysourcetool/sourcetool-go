package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type testServer struct {
	*httptest.Server
	connCh chan *websocket.Conn
}

func newTestServer() *testServer {
	connCh := make(chan *websocket.Conn, 1)
	s := &testServer{
		connCh: connCh,
	}

	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// APIキーの検証
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		s.connCh <- conn
	}))

	return s
}

func (s *testServer) WaitForConnection(t *testing.T) *websocket.Conn {
	t.Helper()
	select {
	case conn := <-s.connCh:
		return conn
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for connection")
		return nil
	}
}

func TestNewClient(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	client, err := NewClient(Config{
		URL:    wsURL,
		APIKey: "test_api_key",
	})

	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	conn := s.WaitForConnection(t)
	if conn == nil {
		t.Fatal("connection not established")
	}
}

func TestClient_MessageHandling(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	client, err := NewClient(Config{
		URL:    wsURL,
		APIKey: "test_api_key",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	conn := s.WaitForConnection(t)
	if conn == nil {
		t.Fatal("connection not established")
	}

	// テスト用のメッセージハンドラー
	receivedCh := make(chan *Message, 1)
	client.RegisterHandler(MessageMethodInitializeClient, func(msg *Message) error {
		receivedCh <- msg
		return nil
	})

	// テストメッセージの送信
	testMsg := &Message{
		ID:     "test_id",
		Kind:   MessageKindCall,
		Method: MessageMethodInitializeClient,
		Payload: json.RawMessage(`{
			"sessionId": "test_session",
			"pageId": "test_page"
		}`),
	}

	if err := conn.WriteJSON(testMsg); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	// メッセージの受信を待つ
	select {
	case msg := <-receivedCh:
		if msg.ID != testMsg.ID {
			t.Errorf("message ID = %v, want %v", msg.ID, testMsg.ID)
		}
		if msg.Method != testMsg.Method {
			t.Errorf("message Method = %v, want %v", msg.Method, testMsg.Method)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestClient_EnqueueWithResponse(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	client, err := NewClient(Config{
		URL:    wsURL,
		APIKey: "test_api_key",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	conn := s.WaitForConnection(t)
	if conn == nil {
		t.Fatal("connection not established")
	}

	// レスポンスを待つgoroutineを開始
	go func() {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			t.Errorf("failed to read message: %v", err)
			return
		}

		// レスポンスを送信
		resp := Message{
			ID:      msg.ID,
			Kind:    MessageKindResponse,
			Method:  msg.Method,
			Payload: json.RawMessage(`{"status":"ok"}`),
		}
		if err := conn.WriteJSON(&resp); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}()

	// メッセージを送信して応答を待つ
	resp, err := client.EnqueueWithResponse("test_id", MessageMethodInitializeHost, map[string]string{
		"test": "data",
	})
	if err != nil {
		t.Fatalf("failed to get response: %v", err)
	}

	if resp.Kind != MessageKindResponse {
		t.Errorf("response Kind = %v, want %v", resp.Kind, MessageKindResponse)
	}
}

func TestClient_Reconnection(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	reconnected := make(chan struct{})
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	client, err := NewClient(Config{
		URL:            wsURL,
		APIKey:         "test_api_key",
		PingInterval:   100 * time.Millisecond,
		ReconnectDelay: 100 * time.Millisecond,
		OnReconnected: func() {
			reconnected <- struct{}{}
		},
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// 最初の接続を確認
	conn := s.WaitForConnection(t)
	if conn == nil {
		t.Fatal("initial connection not established")
	}

	// 接続を切断
	conn.Close()

	// 再接続を待つ
	select {
	case <-reconnected:
		// 再接続成功
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for reconnection")
	}

	// 新しい接続を確認
	conn = s.WaitForConnection(t)
	if conn == nil {
		t.Fatal("reconnection not established")
	}
}
