package textinput

type Options struct {
	Label        string
	Placeholder  string
	DefaultValue string
	Required     bool
	MaxLength    *int
	MinLength    *int
}

type Option func(*Options)
