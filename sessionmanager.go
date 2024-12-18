package sourcetool

import (
	"sync"

	"github.com/gofrs/uuid/v5"
)

type SessionManager struct {
	activeSessions map[uuid.UUID]*Session
	// TODO: Manage sessions that were not explicitly CLOSE_SESSION
	disconnectedSessions map[uuid.UUID]*Session
	mu                   sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		activeSessions:       make(map[uuid.UUID]*Session),
		disconnectedSessions: make(map[uuid.UUID]*Session),
	}
}

func (s *SessionManager) SetSession(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeSessions[session.ID] = session
}

func (s *SessionManager) GetSession(id uuid.UUID) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeSessions[id]
}
