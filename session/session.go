package session

import "github.com/gofrs/uuid/v5"

type Session struct {
	ID    uuid.UUID
	State *State
}
