package dateinput

import "time"

type Options struct {
	Label        string
	Placeholder  string
	DefaultValue *time.Time
	Required     bool
	Disabled     bool
	Format       string
	MaxValue     *time.Time
	MinValue     *time.Time
	Location     *time.Location
}

type Option func(*Options)
