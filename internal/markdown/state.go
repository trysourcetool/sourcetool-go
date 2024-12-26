package markdown

import "github.com/gofrs/uuid/v5"

type State struct {
	ID   uuid.UUID `json:"-"`
	Body string    `json:"body"`
}
