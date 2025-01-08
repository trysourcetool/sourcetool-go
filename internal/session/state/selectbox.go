package state

import "github.com/gofrs/uuid/v5"

const WidgetTypeSelectbox WidgetType = "selectbox"

type SelectboxState struct {
	ID           uuid.UUID
	Label        string
	Value        *int
	Options      []string
	Placeholder  string
	DefaultValue *int
	Required     bool
	Disabled     bool
}

func (s *SelectboxState) IsWidgetState()      {}
func (s *SelectboxState) GetType() WidgetType { return WidgetTypeSelectbox }
