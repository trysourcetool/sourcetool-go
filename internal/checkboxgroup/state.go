package checkboxgroup

import "github.com/gofrs/uuid/v5"

const WidgetType = "checkboxGroup"

type Value struct {
	Values  []string
	Indexes []int
}

type State struct {
	ID           uuid.UUID
	Label        string
	Value        []int
	Options      []string
	DefaultValue []int
	Required     bool
	Disabled     bool
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
