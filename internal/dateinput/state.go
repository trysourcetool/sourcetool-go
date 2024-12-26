package dateinput

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type State struct {
	ID           uuid.UUID  `json:"-"`
	Value        time.Time  `json:"value"`
	Label        string     `json:"label"`
	Placeholder  string     `json:"placeholder"`
	DefaultValue *time.Time `json:"defaultValue"`
	Required     bool       `json:"required"`
	Format       string     `json:"format"`
	MaxValue     *time.Time `json:"maxValue"`
	MinValue     *time.Time `json:"minValue"`
}
