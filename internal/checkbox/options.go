package checkbox

type Options struct {
	Label        string
	DefaultValue bool
	Required     bool
}

type Option func(*Options)
