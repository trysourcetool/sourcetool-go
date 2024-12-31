package checkbox

type Options struct {
	Label        string
	DefaultValue bool
	Required     bool
	Disabled     bool
}

type Option func(*Options)
