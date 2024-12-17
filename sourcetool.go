package sourcetool

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gofrs/uuid/v5"
)

type Sourcetool struct {
	apiKey      string
	endpoint    string
	subdomain   string
	navigations []*Navigation
	pages       map[uuid.UUID]*Page
	mu          sync.RWMutex
}

func New(apiKey string) *Sourcetool {
	subdomain := strings.Split(apiKey, "_")[0]
	s := &Sourcetool{
		apiKey:      apiKey,
		subdomain:   subdomain,
		endpoint:    fmt.Sprintf("ws://%s.localhost:8080/ws", subdomain),
		navigations: make([]*Navigation, 0),
		pages:       make(map[uuid.UUID]*Page),
	}
	return s
}

func (s *Sourcetool) Listen() error {
	StartRuntime(s.apiKey, s.endpoint, s.pages)
	defer Runtime.CloseConnection()

	return Runtime.Wait()
}

func (s *Sourcetool) Close() error {
	return Runtime.CloseConnection()
}
