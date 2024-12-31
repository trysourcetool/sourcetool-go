package timeinput

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const WidgetType = "timeInput"

type State struct {
	ID           uuid.UUID
	Value        *time.Time
	Label        string
	Placeholder  string
	DefaultValue *time.Time
	Required     bool
	Disabled     bool
	Location     *time.Location
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
