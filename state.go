package sourcetool

import (
	"sync"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/textinput"
)

type State struct {
	data map[uuid.UUID]any // ui ID -> options state
	mu   sync.RWMutex
}

func NewState() *State {
	return &State{
		data: make(map[uuid.UUID]any),
	}
}

// GetState returns the state for the given UI component ID
func (s *State) Get(id uuid.UUID) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[id]
}

func (s *State) GetTextInput(id uuid.UUID) *textinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.data[id].(*textinput.State)
	if !ok {
		return nil
	}
	return state
}

// SetState sets the state for the given UI component ID
func (s *State) Set(id uuid.UUID, state any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = state
}
