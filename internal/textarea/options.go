package textarea

type Options struct {
	Label        string
	Placeholder  string
	DefaultValue string
	Required     bool
	MaxLength    *int
	MinLength    *int
	MaxLines     *int
	MinLines     *int
	AutoResize   bool
}

type Option func(*Options)
