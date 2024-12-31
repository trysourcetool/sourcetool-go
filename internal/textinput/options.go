package textinput

type Options struct {
	Label        string
	Placeholder  string
	DefaultValue string
	Required     bool
	Disabled     bool
	MaxLength    *int
	MinLength    *int
}

type Option func(*Options)
