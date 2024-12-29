package multiselect

import "github.com/trysourcetool/sourcetool-go/internal/multiselect"

func Options(options ...string) multiselect.Option {
	return func(o *multiselect.Options) {
		o.Options = options
	}
}

func Placeholder(placeholder string) multiselect.Option {
	return func(o *multiselect.Options) {
		o.Placeholder = placeholder
	}
}

func DefaultValue(defaultValue ...string) multiselect.Option {
	return func(o *multiselect.Options) {
		o.DefaultValue = defaultValue
	}
}

func Required(required bool) multiselect.Option {
	return func(o *multiselect.Options) {
		o.Required = required
	}
}

func FormatFunc(formatFunc func(string, int) string) multiselect.Option {
	return func(o *multiselect.Options) {
		o.FormatFunc = formatFunc
	}
}
