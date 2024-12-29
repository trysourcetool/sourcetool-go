package multiselect

import "github.com/gofrs/uuid/v5"

const WidgetType = "multiSelect"

type Value struct {
	Values  []string
	Indexes []int
}

type State struct {
	ID           uuid.UUID
	Label        string
	Value        *Value
	Options      []string
	Placeholder  string
	DefaultValue []string
	Required     bool
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
