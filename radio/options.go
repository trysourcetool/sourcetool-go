package radio

import "github.com/trysourcetool/sourcetool-go/internal/radio"

func Options(options ...string) radio.Option {
	return func(o *radio.Options) {
		o.Options = options
	}
}

func DefaultValue(defaultValue string) radio.Option {
	return func(o *radio.Options) {
		o.DefaultValue = &defaultValue
	}
}

func Required(required bool) radio.Option {
	return func(o *radio.Options) {
		o.Required = required
	}
}

func Disabled(disabled bool) radio.Option {
	return func(o *radio.Options) {
		o.Disabled = disabled
	}
}

func FormatFunc(formatFunc func(string, int) string) radio.Option {
	return func(o *radio.Options) {
		o.FormatFunc = formatFunc
	}
}
