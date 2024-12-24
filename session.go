package sourcetool

import (
	"sync"

	"github.com/gofrs/uuid/v5"
	"github.com/trysourcetool/sourcetool-go/button"
	"github.com/trysourcetool/sourcetool-go/table"
	"github.com/trysourcetool/sourcetool-go/textinput"
)

type session struct {
	id    uuid.UUID
	state *state
}

func newSession(id uuid.UUID) *session {
	return &session{
		id:    id,
		state: newState(),
	}
}

type sessionManager struct {
	activeSessions map[uuid.UUID]*session
	// TODO: Manage sessions that were not explicitly CLOSE_SESSION
	disconnectedSessions map[uuid.UUID]*session
	mu                   sync.RWMutex
}

func newSessionManager() *sessionManager {
	return &sessionManager{
		activeSessions:       make(map[uuid.UUID]*session),
		disconnectedSessions: make(map[uuid.UUID]*session),
	}
}

func (s *sessionManager) getSession(id uuid.UUID) *session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeSessions[id]
}

func (s *sessionManager) setSession(session *session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeSessions[session.id] = session
}

func (s *sessionManager) deleteSession(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.activeSessions, id)
}

type state struct {
	data map[uuid.UUID]any // ui ID -> options state
	mu   sync.RWMutex
}

func newState() *state {
	return &state{
		data: make(map[uuid.UUID]any),
	}
}

func (s *state) getTextInput(id uuid.UUID) *textinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.data[id].(*textinput.State)
	if !ok {
		return nil
	}
	return state
}

func (s *state) getTable(id uuid.UUID) *table.State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.data[id].(*table.State)
	if !ok {
		return nil
	}
	return state
}

func (s *state) getButton(id uuid.UUID) *button.State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.data[id].(*button.State)
	if !ok {
		return nil
	}
	return state
}

func (s *state) resetButtons() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, state := range s.data {
		if buttonState, ok := state.(*button.State); ok {
			buttonState.Value = false
		}
	}
}

func (s *state) set(id uuid.UUID, state any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = state
}

func (s *state) setStates(states map[uuid.UUID]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, state := range states {
		s.data[id] = state
	}
}
