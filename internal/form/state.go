package form

import "github.com/gofrs/uuid/v5"

type State struct {
	ID            uuid.UUID `json:"-"`
	ClearOnSubmit bool      `json:"clearOnSubmit"`
}
