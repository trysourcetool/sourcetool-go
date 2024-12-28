package columns

import "github.com/gofrs/uuid/v5"

const WidgetType = "columns"

type State struct {
	ID      uuid.UUID
	Columns int
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
