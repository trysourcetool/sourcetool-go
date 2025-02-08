package session

import (
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
)

const (
	// TTL for disconnected sessions
	disconnectedSessionTTL = 2 * time.Minute
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
	activeSessions       map[uuid.UUID]*Session
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

	if disconnectedSession, ok := s.disconnectedSessions[session.ID]; ok {
		session.State = disconnectedSession.State
		delete(s.disconnectedSessions, session.ID)
	}

	s.activeSessions[session.ID] = session
}

func (s *SessionManager) DisconnectSession(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if session, ok := s.activeSessions[id]; ok {
		s.disconnectedSessions[id] = session
		delete(s.activeSessions, id)

		go func(sessionID uuid.UUID) {
			time.Sleep(disconnectedSessionTTL)
			s.mu.Lock()
			delete(s.disconnectedSessions, sessionID)
			s.mu.Unlock()
		}(id)
	}
}
