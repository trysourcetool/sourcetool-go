package dateinput

import "time"

type Options struct {
	Label        string
	Placeholder  string
	DefaultValue *time.Time
	Required     bool
	Format       string
	MaxValue     *time.Time
	MinValue     *time.Time
}

type Option func(*Options)
