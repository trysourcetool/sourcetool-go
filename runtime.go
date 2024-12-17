package sourcetool

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"

	ws "github.com/trysourcetool/sourcetool-go/websocket"
)

var Runtime *runtime
var once sync.Once

type runtime struct {
	wsClient       ws.Client
	sessionManager *SessionManager
	pageManager    *PageManager
}

func StartRuntime(apiKey, endpoint string, pages map[uuid.UUID]*Page) {
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
				initializeHost(apiKey, pages)
			},
		})
		if err != nil {
			log.Fatalf("failed to create websocket client: %v", err)
		}

		Runtime = &runtime{
			wsClient:       wsClient,
			sessionManager: NewSessionManager(),
			pageManager:    NewPageManager(pages),
		}

		initializeHost(apiKey, pages)
	})
}

func initializeHost(apiKey string, pages map[uuid.UUID]*Page) {
	pagesPayload := make([]*ws.InitializeHostPagePayload, 0, len(pages))
	for _, page := range pages {
		pagesPayload = append(pagesPayload, &ws.InitializeHostPagePayload{
			ID:   page.ID.String(),
			Name: page.Name,
		})
	}

	resp, err := Runtime.EnqueueMessageWithResponse(uuid.Must(uuid.NewV4()).String(), ws.MessageMethodInitializeHost, ws.InitializeHostPayload{
		APIKey:     apiKey,
		SDKName:    "sourcetool-go",
		SDKVersion: "0.1.0",
		Pages:      pagesPayload,
	})
	if err != nil {
		log.Fatalf("failed to send initialize host message: %v", err)
	}
	if resp.Error != nil {
		log.Fatalf("initialize host message failed: %v", resp.Error)
	}

	log.Printf("initialize host message sent: %v", resp)
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

func (r *runtime) SetSession(s *Session) {
	r.sessionManager.SetSession(s)
}

func (r *runtime) GetSession(ctx context.Context) (*Session, error) {
	sessionID, err := SessionIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.sessionManager.GetSession(sessionID), nil
}

func (r *runtime) GetPage(id uuid.UUID) *Page {
	return r.pageManager.GetPage(id)
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

		sessionID, err := uuid.FromString(payload.SessionID)
		if err != nil {
			return err
		}
		pageID, err := uuid.FromString(payload.PageID)
		if err != nil {
			return err
		}

		log.Println("Creating new session with ID:", sessionID)

		// Create new session with state
		Runtime.SetSession(NewSession(sessionID))

		// Create context with session ID for this connection
		// ctx := session.WithSessionID(context.Background(), sessionID)

		page := Runtime.pageManager.GetPage(pageID)
		if page == nil {
			return fmt.Errorf("page not found: %s", pageID)
		}

		ctx := &Context{
			session: NewSession(sessionID),
			page:    page,
			cursor:  NewCursor(MAIN),
		}

		// Run pages with session context
		if err := page.Run(ctx); err != nil {
			return fmt.Errorf("failed to run pages: %v", err)
		}

		return nil

	case ws.MessageMethodCloseSession:
		log.Println("Received CloseSession message")
		var payload ws.CloseSessionPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %v", err)
		}

		// Session cleanup can be handled by the session manager
		// No need to track current session as it's in the context
		return nil

	default:
		return fmt.Errorf("unknown message method: %s", msg.Method)
	}
}
