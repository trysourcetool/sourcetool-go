package numberinput

import "github.com/trysourcetool/sourcetool-go/internal/numberinput"

func Placeholder(placeholder string) numberinput.Option {
	return func(opts *numberinput.Options) {
		opts.Placeholder = placeholder
	}
}

func DefaultValue(value float64) numberinput.Option {
	return func(opts *numberinput.Options) {
		opts.DefaultValue = &value
	}
}

func Required(required bool) numberinput.Option {
	return func(opts *numberinput.Options) {
		opts.Required = required
	}
}

func Disabled(disabled bool) numberinput.Option {
	return func(opts *numberinput.Options) {
		opts.Disabled = disabled
	}
}

func MaxValue(value float64) numberinput.Option {
	return func(opts *numberinput.Options) {
		opts.MaxValue = &value
	}
}

func MinValue(value float64) numberinput.Option {
	return func(opts *numberinput.Options) {
		opts.MinValue = &value
	}
}
