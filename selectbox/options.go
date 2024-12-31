package selectbox

import "github.com/trysourcetool/sourcetool-go/internal/selectbox"

func Options(options ...string) selectbox.Option {
	return func(o *selectbox.Options) {
		o.Options = options
	}
}

func Placeholder(placeholder string) selectbox.Option {
	return func(o *selectbox.Options) {
		o.Placeholder = placeholder
	}
}

func DefaultValue(defaultValue string) selectbox.Option {
	return func(o *selectbox.Options) {
		o.DefaultValue = &defaultValue
	}
}

func Required(required bool) selectbox.Option {
	return func(o *selectbox.Options) {
		o.Required = required
	}
}

func Disabled(disabled bool) selectbox.Option {
	return func(o *selectbox.Options) {
		o.Disabled = disabled
	}
}

func FormatFunc(formatFunc func(string, int) string) selectbox.Option {
	return func(o *selectbox.Options) {
		o.FormatFunc = formatFunc
	}
}
