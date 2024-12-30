package datetimeinput

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const WidgetType = "dateTimeInput"

type State struct {
	ID           uuid.UUID
	Value        *time.Time
	Label        string
	Placeholder  string
	DefaultValue *time.Time
	Required     bool
	Format       string
	MaxValue     *time.Time
	MinValue     *time.Time
	Location     *time.Location
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
