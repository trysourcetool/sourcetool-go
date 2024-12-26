package session

import (
	"sync"

	"github.com/gofrs/uuid/v5"
)

type Session struct {
	ID     uuid.UUID
	PageID uuid.UUID
	State  *State
}

func New(id uuid.UUID, pageID uuid.UUID) *Session {
	return &Session{
		ID:     id,
		PageID: pageID,
		State:  newState(),
	}
}

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

func (s *SessionManager) GetSession(id uuid.UUID) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeSessions[id]
}

func (s *SessionManager) SetSession(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeSessions[session.ID] = session
}

func (s *SessionManager) DeleteSession(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.activeSessions, id)
}
