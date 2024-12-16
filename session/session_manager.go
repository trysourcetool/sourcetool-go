package session

import "github.com/gofrs/uuid/v5"

type SessionManager struct {
	activeSessions map[uuid.UUID]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		activeSessions: make(map[uuid.UUID]*Session),
	}
}

func (s *SessionManager) SetSession(session *Session) {
	s.activeSessions[session.ID] = session
}

func (s *SessionManager) GetSession(id uuid.UUID) *Session {
	return s.activeSessions[id]
}
