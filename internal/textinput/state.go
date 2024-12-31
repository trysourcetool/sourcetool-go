package textinput

import "github.com/gofrs/uuid/v5"

const WidgetType = "textInput"

type State struct {
	ID           uuid.UUID
	Value        string
	Label        string
	Placeholder  string
	DefaultValue string
	Required     bool
	Disabled     bool
	MaxLength    *int
	MinLength    *int
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
