package columnitem

import "github.com/gofrs/uuid/v5"

type State struct {
	ID     uuid.UUID `json:"-"`
	Weight float64   `json:"weight"`
}
