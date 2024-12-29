package selectbox

import "github.com/gofrs/uuid/v5"

const WidgetType = "selectbox"

type State struct {
	ID           uuid.UUID
	Label        string
	Value        *any
	Values       []any
	Labels       []string
	Placeholder  string
	DefaultIndex *int
	Required     bool
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
