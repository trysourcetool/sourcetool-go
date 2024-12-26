package numberinput

import "github.com/gofrs/uuid/v5"

type ReturnValue float64

type State struct {
	ID           uuid.UUID   `json:"-"`
	Value        ReturnValue `json:"value"`
	Label        string      `json:"label"`
	Placeholder  string      `json:"placeholder"`
	DefaultValue float64     `json:"defaultValue"`
	Required     bool        `json:"required"`
	MaxValue     *float64    `json:"maxValue"`
	MinValue     *float64    `json:"minValue"`
}
