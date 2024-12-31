package textarea

type Options struct {
	Label        string
	Placeholder  string
	DefaultValue string
	Required     bool
	Disabled     bool
	MaxLength    *int
	MinLength    *int
	MaxLines     *int
	MinLines     *int
	AutoResize   bool
}

type Option func(*Options)
