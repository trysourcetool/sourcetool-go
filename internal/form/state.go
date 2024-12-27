package form

import "github.com/gofrs/uuid/v5"

const WidgetType = "form"

type State struct {
	ID             uuid.UUID `json:"-"`
	Value          bool      `json:"value"`
	ButtonLabel    string    `json:"buttonLabel"`
	ButtonDisabled bool      `json:"buttonDisabled"`
	ClearOnSubmit  bool      `json:"clearOnSubmit"`
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
