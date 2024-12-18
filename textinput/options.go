package textinput

type Options struct {
	Label        string
	Placeholder  string
	DefaultValue string
	Required     bool
	MaxLength    int
	MinLength    int
}

type Option func(*Options)

func WithPlaceholder(placeholder string) Option {
	return func(opts *Options) {
		opts.Placeholder = placeholder
	}
}

func WithDefaultValue(value string) Option {
	return func(opts *Options) {
		opts.DefaultValue = value
	}
}

func WithRequired(required bool) Option {
	return func(opts *Options) {
		opts.Required = required
	}
}

func WithMaxLength(length int) Option {
	return func(opts *Options) {
		opts.MaxLength = length
	}
}

func WithMinLength(length int) Option {
	return func(opts *Options) {
		opts.MinLength = length
	}
}
