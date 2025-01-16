package sourcetool

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gofrs/uuid/v5"
)

type Sourcetool struct {
	*pageBuilder

	apiKey    string
	endpoint  string
	subdomain string
	runtime   *runtime
	pages     map[uuid.UUID]*page
	mu        sync.RWMutex
}

func New(apiKey string) *Sourcetool {
	subdomain := strings.Split(apiKey, "_")[0]
	s := &Sourcetool{
		apiKey:    apiKey,
		subdomain: subdomain,
		endpoint:  fmt.Sprintf("ws://%s.local.trysourcetool.com:8080/ws", subdomain),
		pages:     make(map[uuid.UUID]*page),
		pageBuilder: &pageBuilder{
			namespaceDNS: fmt.Sprintf("%s.trysourcetool.com", subdomain),
		},
	}
	s.pageBuilder.sourcetool = s
	return s
}

func (s *Sourcetool) Listen() error {
	if err := s.validatePages(); err != nil {
		return err
	}

	r, err := startRuntime(s.apiKey, s.endpoint, s.pages)
	if err != nil {
		return err
	}
	defer r.wsClient.Close()

	s.runtime = r

	return r.wsClient.Wait()
}

func (s *Sourcetool) Close() error {
	return s.runtime.wsClient.Close()
}

func (s *Sourcetool) validatePages() error {
	pageNames := make(map[string]struct{})
	for _, p := range s.pages {
		if p.name == "" {
			return errors.New("page name cannot be empty")
		}
		if _, exists := pageNames[p.name]; exists {
			return fmt.Errorf("duplicate page name: %s", p.name)
		}
		pageNames[p.name] = struct{}{}
	}
	return nil
}
