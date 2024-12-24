package sourcetool

import (
	"log"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"

	ws "github.com/trysourcetool/sourcetool-go/websocket"
)

var once sync.Once

type runtime struct {
	wsClient       ws.Client
	sessionManager *sessionManager
	pageManager    *pageManager
}

func startRuntime(apiKey, endpoint string, pages map[uuid.UUID]*page) *runtime {
	var r *runtime
	once.Do(func() {
		r = &runtime{
			sessionManager: newSessionManager(),
			pageManager:    newPageManager(pages),
		}

		wsClient, err := ws.NewClient(ws.Config{
			URL:            endpoint,
			APIKey:         apiKey,
			PingInterval:   1 * time.Second,
			ReconnectDelay: 1 * time.Second,
			OnReconnecting: func() {
				log.Println("Reconnecting...")
			},
			OnReconnected: func() {
				log.Println("Reconnected!")
				r.initializeHost(apiKey, pages)
			},
		})
		if err != nil {
			log.Fatalf("failed to create websocket client: %v", err)
		}

		r.wsClient = wsClient
		msgHandler := &messageHandler{r}
		wsClient.RegisterHandler(ws.MessageMethodInitializeClient, msgHandler.initializeCilent)
		wsClient.RegisterHandler(ws.MessageMethodRerunPage, msgHandler.rerunPage)
		wsClient.RegisterHandler(ws.MessageMethodCloseSession, msgHandler.closeSession)

		r.initializeHost(apiKey, pages)
	})

	return r
}

func (r *runtime) initializeHost(apiKey string, pages map[uuid.UUID]*page) {
	pagesPayload := make([]*ws.InitializeHostPagePayload, 0, len(pages))
	for _, page := range pages {
		pagesPayload = append(pagesPayload, &ws.InitializeHostPagePayload{
			ID:   page.id.String(),
			Name: page.name,
		})
	}

	resp, err := r.wsClient.EnqueueWithResponse(uuid.Must(uuid.NewV4()).String(), ws.MessageMethodInitializeHost, ws.InitializeHostPayload{
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
