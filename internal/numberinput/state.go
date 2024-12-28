package numberinput

import "github.com/gofrs/uuid/v5"

const WidgetType = "numberInput"

type State struct {
	ID           uuid.UUID
	Value        *float64
	Label        string
	Placeholder  string
	DefaultValue *float64
	Required     bool
	MaxValue     *float64
	MinValue     *float64
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
