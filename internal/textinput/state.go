package textinput

import "github.com/gofrs/uuid/v5"

const WidgetType = "textInput"

type State struct {
	ID           uuid.UUID `json:"-"`
	Value        string    `json:"value"`
	Label        string    `json:"label"`
	Placeholder  string    `json:"placeholder"`
	DefaultValue string    `json:"defaultValue"`
	Required     bool      `json:"required"`
	MaxLength    *int      `json:"maxLength"`
	MinLength    *int      `json:"minLength"`
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
