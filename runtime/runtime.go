package runtime

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/session"
	ws "github.com/trysourcetool/sourcetool-go/websocket"
)

var Runtime *runtime
var once sync.Once

type runtime struct {
	wsClient       ws.Client
	sessionManager *session.SessionManager
}

func Start(apiKey, endpoint string) {
	once.Do(func() {
		wsClient, err := ws.NewClient(ws.Config{
			URL:            endpoint,
			APIKey:         apiKey,
			PingInterval:   1 * time.Second,
			ReconnectDelay: 1 * time.Second,
			CallHandler:    callHandler,
			OnReconnecting: func() {
				log.Println("Reconnecting...")
			},
			OnReconnected: func() {
				log.Println("Reconnected!")
			},
		})
		if err != nil {
			log.Fatalf("failed to create websocket client: %v", err)
		}

		Runtime = &runtime{
			wsClient:       wsClient,
			sessionManager: session.NewSessionManager(),
		}
	})
}

func (r *runtime) CloseConnection() error {
	return r.wsClient.Close()
}

func (r *runtime) EnqueueMessage(id string, method ws.MessageMethod, data any) {
	r.wsClient.Enqueue(id, method, data)
}

func (r *runtime) EnqueueMessageWithResponse(id string, method ws.MessageMethod, data any) (*ws.Message, error) {
	resp, err := r.wsClient.EnqueueWithResponse(id, method, data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *runtime) Wait() error {
	return r.wsClient.Wait()
}

func (r *runtime) SetSession(s *session.Session) {
	r.sessionManager.SetSession(s)
}

func callHandler(msg *ws.Message) error {
	var payload map[string]any
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	log.Printf("Received message: %+v", msg)
	log.Printf("Received payload: %+v", payload)

	switch msg.Method {
	case ws.MessageMethodInitializeClient:
		log.Println("Received InitializeClient message")

		var payload ws.InitializeClientPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %v", err)
		}

		// TODO: Set session
		// TODO: Run pages
		sessionID, err := uuid.FromString(payload.SessionID)
		if err != nil {
			return err
		}
		Runtime.SetSession(&session.Session{
			ID: sessionID,
		})
		return nil
	case ws.MessageMethodCloseSession:
		log.Println("Received InitializeClient message")
		// TODO: Delete session
		return nil
	default:
		return fmt.Errorf("unknown message method: %s", msg.Method)
	}
}
