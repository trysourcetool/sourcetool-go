package checkboxgroup

import "github.com/trysourcetool/sourcetool-go/internal/checkboxgroup"

func Options(options ...string) checkboxgroup.Option {
	return func(o *checkboxgroup.Options) {
		o.Options = options
	}
}

func DefaultValue(defaultValue ...string) checkboxgroup.Option {
	return func(o *checkboxgroup.Options) {
		o.DefaultValue = defaultValue
	}
}

func Required(required bool) checkboxgroup.Option {
	return func(o *checkboxgroup.Options) {
		o.Required = required
	}
}

func Disabled(disabled bool) checkboxgroup.Option {
	return func(o *checkboxgroup.Options) {
		o.Disabled = disabled
	}
}

func FormatFunc(formatFunc func(string, int) string) checkboxgroup.Option {
	return func(o *checkboxgroup.Options) {
		o.FormatFunc = formatFunc
	}
}
