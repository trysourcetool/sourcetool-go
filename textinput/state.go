package textinput

import "github.com/gofrs/uuid/v5"

type State struct {
	ID           uuid.UUID `json:"-"`
	Value        string    `json:"value"`
	Label        string    `json:"label"`
	Placeholder  string    `json:"placeholder"`
	DefaultValue string    `json:"defaultValue"`
	Required     bool      `json:"required"`
	MaxLength    int       `json:"maxLength"`
	MinLength    int       `json:"minLength"`
}
