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
	runtime     *runtime
	navigations []*navigation
	pages       map[uuid.UUID]*page
	mu          sync.RWMutex
}

func New(apiKey string) *Sourcetool {
	subdomain := strings.Split(apiKey, "_")[0]
	s := &Sourcetool{
		apiKey:      apiKey,
		subdomain:   subdomain,
		endpoint:    fmt.Sprintf("ws://%s.local.trysourcetool.com:8080/ws", subdomain),
		navigations: make([]*navigation, 0),
		pages:       make(map[uuid.UUID]*page),
	}
	return s
}

func (s *Sourcetool) Listen() error {
	r := startRuntime(s.apiKey, s.endpoint, s.pages)
	defer r.wsClient.Close()

	s.runtime = r

	return r.wsClient.Wait()
}

func (s *Sourcetool) Close() error {
	return s.runtime.wsClient.Close()
}
