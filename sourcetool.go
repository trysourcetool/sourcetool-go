package sourcetool

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/runtime"
	"github.com/trysourcetool/sourcetool-go/ui"
	ws "github.com/trysourcetool/sourcetool-go/websocket"
)

type Sourcetool struct {
	apiKey    string
	endpoint  string
	subdomain string
	pages     []*ui.Page
	mu        sync.Mutex
}

func New(apiKey string) *Sourcetool {
	subdomain := strings.Split(apiKey, "_")[0]
	s := &Sourcetool{
		apiKey:    apiKey,
		subdomain: subdomain,
		endpoint:  fmt.Sprintf("ws://%s.localhost:8080/ws", subdomain),
	}
	return s
}

func (s *Sourcetool) Listen() error {
	runtime.Start(s.apiKey, s.endpoint)
	defer runtime.Runtime.CloseConnection()

	s.initializeHost()

	return runtime.Runtime.Wait()
}

func (s *Sourcetool) RegisterPages(pages ...*ui.Page) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, page := range pages {
		page.ID = s.generatePageID(page.Name)
		s.pages = append(s.pages, page)
	}
}

func (s *Sourcetool) initializeHost() {
	pages := make([]*ws.InitializeHostPagePayload, 0, len(s.pages))
	for _, page := range s.pages {
		pages = append(pages, &ws.InitializeHostPagePayload{
			ID:   page.ID,
			Name: page.Name,
		})
	}

	resp, err := runtime.Runtime.EnqueueMessageWithResponse(uuid.Must(uuid.NewV4()).String(), ws.MessageMethodInitializeHost, ws.InitializeHostPayload{
		APIKey:     s.apiKey,
		SDKName:    "sourcetool-go",
		SDKVersion: "0.1.0",
		Pages:      pages,
	})
	if err != nil {
		log.Fatalf("failed to send initialize host message: %v", err)
	}
	if resp.Error != nil {
		log.Fatalf("initialize host message failed: %v", resp.Error)
	}

	log.Printf("initialize host message sent: %v", resp)
}

func (s *Sourcetool) generatePageID(pageName string) string {
	ns := uuid.NewV5(uuid.NamespaceDNS, fmt.Sprintf("%s.trysourcetool.com", s.subdomain))
	return uuid.NewV5(ns, pageName).String()
}

func (s *Sourcetool) Close() error {
	return runtime.Runtime.CloseConnection()
}
