package button

import "github.com/gofrs/uuid/v5"

type ReturnValue bool

type State struct {
	ID       uuid.UUID   `json:"-"`
	Value    ReturnValue `json:"value"`
	Label    string      `json:"label"`
	Disabled bool        `json:"disabled"`
}
