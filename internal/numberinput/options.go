package numberinput

type Options struct {
	Label        string
	Placeholder  string
	DefaultValue *float64
	Required     bool
	Disabled     bool
	MaxValue     *float64
	MinValue     *float64
}

type Option func(*Options)
