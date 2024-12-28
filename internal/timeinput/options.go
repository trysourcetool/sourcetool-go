package timeinput

import "time"

type Options struct {
	Label        string
	Placeholder  string
	DefaultValue *time.Time
	Required     bool
	Location     *time.Location
}

type Option func(*Options)
