package textinput

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.TextInputOptions)
}

type placeholderOption string

func (p placeholderOption) Apply(opts *options.TextInputOptions) {
	opts.Placeholder = string(p)
}

func Placeholder(placeholder string) Option {
	return placeholderOption(placeholder)
}

type defaultValueOption string

func (d defaultValueOption) Apply(opts *options.TextInputOptions) {
	opts.DefaultValue = string(d)
}

func DefaultValue(value string) Option {
	return defaultValueOption(value)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.TextInputOptions) {
	opts.Required = bool(r)
}

func Required(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.TextInputOptions) {
	opts.Disabled = bool(d)
}

func Disabled(disabled bool) Option {
	return disabledOption(disabled)
}

type maxLengthOption int

func (m maxLengthOption) Apply(opts *options.TextInputOptions) {
	opts.MaxLength = (*int)(&m)
}

func MaxLength(length int) Option {
	return maxLengthOption(length)
}

type minLengthOption int

func (m minLengthOption) Apply(opts *options.TextInputOptions) {
	opts.MinLength = (*int)(&m)
}

func MinLength(length int) Option {
	return minLengthOption(length)
}
