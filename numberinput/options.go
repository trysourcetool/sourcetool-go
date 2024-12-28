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

func MaxValue(value float64) numberinput.Option {
	return func(opts *numberinput.Options) {
		opts.MinValue = &value
	}
}

func MinLength(value float64) numberinput.Option {
	return func(opts *numberinput.Options) {
		opts.MaxValue = &value
	}
}
