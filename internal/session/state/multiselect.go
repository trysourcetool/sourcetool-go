package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeMultiSelect WidgetType = "multiSelect"

type MultiSelectState struct {
	ID           uuid.UUID
	Label        string
	Value        []int
	Options      []string
	Placeholder  string
	DefaultValue []int
	Required     bool
	Disabled     bool
}

func (s *MultiSelectState) IsWidgetState()      {}
func (s *MultiSelectState) GetType() WidgetType { return WidgetTypeMultiSelect }
