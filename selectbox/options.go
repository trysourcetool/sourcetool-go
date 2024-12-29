package selectbox

import "github.com/trysourcetool/sourcetool-go/internal/selectbox"

func Options(options ...any) selectbox.Option {
	return func(o *selectbox.Options) {
		o.Options = options
	}
}

func Placeholder(placeholder string) selectbox.Option {
	return func(o *selectbox.Options) {
		o.Placeholder = placeholder
	}
}

func DefaultIndex(defaultIndex int) selectbox.Option {
	return func(o *selectbox.Options) {
		o.DefaultIndex = &defaultIndex
	}
}

func Required(required bool) selectbox.Option {
	return func(o *selectbox.Options) {
		o.Required = required
	}
}

func DisplayFunc(displayFunc func(any, int) string) selectbox.Option {
	return func(o *selectbox.Options) {
		o.DisplayFunc = displayFunc
	}
}
