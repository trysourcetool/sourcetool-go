package columnitem

import "github.com/gofrs/uuid/v5"

const WidgetType = "columnItem"

type State struct {
	ID     uuid.UUID
	Weight float64
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
