package button

import "github.com/gofrs/uuid/v5"

const WidgetType = "button"

type State struct {
	ID       uuid.UUID `json:"-"`
	Value    bool      `json:"value"`
	Label    string    `json:"label"`
	Disabled bool      `json:"disabled"`
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
