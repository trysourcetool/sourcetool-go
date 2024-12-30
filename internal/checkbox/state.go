package checkbox

import "github.com/gofrs/uuid/v5"

const WidgetType = "checkbox"

type State struct {
	ID           uuid.UUID
	Label        string
	Value        bool
	DefaultValue bool
	Required     bool
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
