package button

import "github.com/gofrs/uuid/v5"

type State struct {
	ID       uuid.UUID `json:"-"`
	Value    bool      `json:"value"`
	Label    string    `json:"label"`
	Disabled bool      `json:"disabled"`
}
