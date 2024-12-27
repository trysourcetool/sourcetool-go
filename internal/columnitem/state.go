package columnitem

import "github.com/gofrs/uuid/v5"

const WidgetType = "columnItem"

type State struct {
	ID     uuid.UUID `json:"-"`
	Weight float64   `json:"weight"`
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
