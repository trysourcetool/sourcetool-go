package columns

import "github.com/gofrs/uuid/v5"

type State struct {
	ID      uuid.UUID `json:"-"`
	Columns int       `json:"columns"`
}
