package form

import "github.com/gofrs/uuid/v5"

const WidgetType = "form"

type State struct {
	ID             uuid.UUID
	Value          bool
	ButtonLabel    string
	ButtonDisabled bool
	ClearOnSubmit  bool
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
