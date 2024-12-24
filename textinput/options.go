package textinput

import "github.com/trysourcetool/sourcetool-go/internal/textinput"

func Placeholder(placeholder string) textinput.Option {
	return func(opts *textinput.Options) {
		opts.Placeholder = placeholder
	}
}

func DefaultValue(value string) textinput.Option {
	return func(opts *textinput.Options) {
		opts.DefaultValue = value
	}
}

func Required(required bool) textinput.Option {
	return func(opts *textinput.Options) {
		opts.Required = required
	}
}

func MaxLength(length int) textinput.Option {
	return func(opts *textinput.Options) {
		opts.MaxLength = &length
	}
}

func MinLength(length int) textinput.Option {
	return func(opts *textinput.Options) {
		opts.MinLength = &length
	}
}
