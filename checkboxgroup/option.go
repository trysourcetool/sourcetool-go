package checkboxgroup

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.CheckboxGroupOptions)
}

type optionsOption []string

func (o optionsOption) Apply(opts *options.CheckboxGroupOptions) {
	opts.Options = []string(o)
}

func WithOptions(options ...string) Option {
	return optionsOption(options)
}

type defaultValueOption []string

func (d defaultValueOption) Apply(opts *options.CheckboxGroupOptions) {
	opts.DefaultValue = []string(d)
}

func WithDefaultValue(defaultValue ...string) Option {
	return defaultValueOption(defaultValue)
}

type requiredOption bool

func (r requiredOption) Apply(opts *options.CheckboxGroupOptions) {
	opts.Required = bool(r)
}

func WithRequired(required bool) Option {
	return requiredOption(required)
}

type disabledOption bool

func (d disabledOption) Apply(opts *options.CheckboxGroupOptions) {
	opts.Disabled = bool(d)
}

func WithDisabled(disabled bool) Option {
	return disabledOption(disabled)
}

type formatFuncOption func(string, int) string

func (f formatFuncOption) Apply(opts *options.CheckboxGroupOptions) {
	opts.FormatFunc = func(s string, i int) string {
		return f(s, i)
	}
}

func WithFormatFunc(formatFunc func(string, int) string) Option {
	return formatFuncOption(formatFunc)
}
