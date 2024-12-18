package sourcetool

import (
	"github.com/gofrs/uuid/v5"
)

type Session struct {
	ID    uuid.UUID
	State *State
}

func NewSession(id uuid.UUID) *Session {
	return &Session{
		ID:    id,
		State: NewState(),
	}
}
