package sourcetool

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

var once sync.Once

type runtime struct {
	wsClient       websocket.Client
	sessionManager *session.SessionManager
	pageManager    *pageManager
}

func startRuntime(apiKey, endpoint string, pages map[uuid.UUID]*page) *runtime {
	var r *runtime
	once.Do(func() {
		r = &runtime{
			sessionManager: session.NewSessionManager(),
			pageManager:    newPageManager(pages),
		}

		wsClient, err := websocket.NewClient(websocket.Config{
			URL:            endpoint,
			APIKey:         apiKey,
			PingInterval:   1 * time.Second,
			ReconnectDelay: 1 * time.Second,
			OnReconnecting: func() {
				log.Println("Reconnecting...")
			},
			OnReconnected: func() {
				log.Println("Reconnected!")
				r.sendInitializeHost(apiKey, pages)
			},
		})
		if err != nil {
			log.Fatalf("failed to create websocket client: %v", err)
		}

		r.wsClient = wsClient
		wsClient.RegisterHandler(websocket.MessageMethodInitializeClient, r.handleInitializeCilent)
		wsClient.RegisterHandler(websocket.MessageMethodRerunPage, r.handleRerunPage)
		wsClient.RegisterHandler(websocket.MessageMethodCloseSession, r.handleCloseSession)

		r.sendInitializeHost(apiKey, pages)
	})

	return r
}

func (r *runtime) sendInitializeHost(apiKey string, pages map[uuid.UUID]*page) {
	pagesPayload := make([]*websocket.InitializeHostPagePayload, 0, len(pages))
	for _, page := range pages {
		pagesPayload = append(pagesPayload, &websocket.InitializeHostPagePayload{
			ID:   page.id.String(),
			Name: page.name,
		})
	}

	resp, err := r.wsClient.EnqueueWithResponse(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodInitializeHost, websocket.InitializeHostPayload{
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

func (r *runtime) handleInitializeCilent(msg *websocket.Message) error {
	var p websocket.InitializeClientPayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	sessionID, err := uuid.FromString(p.SessionID)
	if err != nil {
		return err
	}
	pageID, err := uuid.FromString(p.PageID)
	if err != nil {
		return err
	}

	log.Println("Creating new session with ID:", sessionID)
	session := session.New(sessionID, pageID)
	r.sessionManager.SetSession(session)

	page := r.pageManager.getPage(pageID)
	if page == nil {
		return fmt.Errorf("page not found: %s", pageID)
	}

	ui := &uiBuilder{
		context: context.Background(),
		runtime: r,
		session: session,
		page:    page,
		cursor:  newCursor(),
	}

	if err := page.run(ui); err != nil {
		return fmt.Errorf("failed to run page: %v", err)
	}

	return nil
}

func (r *runtime) handleRerunPage(msg *websocket.Message) error {
	var p websocket.RerunPagePayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	sessionID, err := uuid.FromString(p.SessionID)
	if err != nil {
		return err
	}
	sess := r.sessionManager.GetSession(sessionID)
	if sess == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	pageID, err := uuid.FromString(p.PageID)
	if err != nil {
		return err
	}
	page := r.pageManager.getPage(pageID)
	if page == nil {
		return fmt.Errorf("page not found: %s", pageID)
	}

	// Reset states if page has changed
	if sess.PageID != pageID {
		sess.State.ResetStates()
	}

	var states map[uuid.UUID]any
	if err := json.Unmarshal(p.State, &states); err != nil {
		return fmt.Errorf("failed to unmarshal state: %v", err)
	}
	sess.State.SetStates(states)

	ui := &uiBuilder{
		context: context.Background(),
		runtime: r,
		session: sess,
		page:    page,
		cursor:  newCursor(),
	}

	if err := page.run(ui); err != nil {
		return fmt.Errorf("failed to run page: %v", err)
	}

	sess.State.ResetButtons()

	return nil
}

func (r *runtime) handleCloseSession(msg *websocket.Message) error {
	var p websocket.CloseSessionPayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	sessionID, err := uuid.FromString(p.SessionID)
	if err != nil {
		return err
	}

	r.sessionManager.DeleteSession(sessionID)

	return nil
}
