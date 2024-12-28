package button

import "github.com/gofrs/uuid/v5"

const WidgetType = "button"

type State struct {
	ID       uuid.UUID
	Value    bool
	Label    string
	Disabled bool
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
