package radio

import "github.com/gofrs/uuid/v5"

const WidgetType = "radio"

type Value struct {
	Value string
	Index int
}

type State struct {
	ID           uuid.UUID
	Label        string
	Value        *int
	Options      []string
	DefaultValue *int
	Required     bool
	Disabled     bool
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
