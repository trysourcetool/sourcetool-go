package textinput

type Options struct {
	Label        string
	Placeholder  string
	DefaultValue string
	Required     bool
	MaxLength    *int
	MinLength    *int
}

func DefaultOptions(label string) *Options {
	return &Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: "",
		Required:     false,
		MaxLength:    nil,
		MinLength:    nil,
	}
}

type Option func(*Options)

func Placeholder(placeholder string) Option {
	return func(opts *Options) {
		opts.Placeholder = placeholder
	}
}

func DefaultValue(value string) Option {
	return func(opts *Options) {
		opts.DefaultValue = value
	}
}

func Required(required bool) Option {
	return func(opts *Options) {
		opts.Required = required
	}
}

func MaxLength(length int) Option {
	return func(opts *Options) {
		opts.MaxLength = &length
	}
}

func MinLength(length int) Option {
	return func(opts *Options) {
		opts.MinLength = &length
	}
}
