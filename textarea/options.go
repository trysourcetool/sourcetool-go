package textarea

import "github.com/trysourcetool/sourcetool-go/internal/textarea"

func Placeholder(placeholder string) textarea.Option {
	return func(opts *textarea.Options) {
		opts.Placeholder = placeholder
	}
}

func DefaultValue(value string) textarea.Option {
	return func(opts *textarea.Options) {
		opts.DefaultValue = value
	}
}

func Required(required bool) textarea.Option {
	return func(opts *textarea.Options) {
		opts.Required = required
	}
}

func MaxLength(length int) textarea.Option {
	return func(opts *textarea.Options) {
		opts.MaxLength = &length
	}
}

func MinLength(length int) textarea.Option {
	return func(opts *textarea.Options) {
		opts.MinLength = &length
	}
}

func MaxLines(lines int) textarea.Option {
	return func(opts *textarea.Options) {
		opts.MaxLines = &lines
	}
}

func MinLines(lines int) textarea.Option {
	return func(opts *textarea.Options) {
		opts.MinLines = &lines
	}
}

func AutoResize(autoResize bool) textarea.Option {
	return func(opts *textarea.Options) {
		opts.AutoResize = autoResize
	}
}
