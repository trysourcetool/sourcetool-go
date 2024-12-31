package textarea

import "github.com/gofrs/uuid/v5"

const WidgetType = "textArea"

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
	MaxLines     *int
	MinLines     *int
	AutoResize   bool
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
