package dateinput

import (
	"time"

	"github.com/trysourcetool/sourcetool-go/internal/dateinput"
)

func Placeholder(placeholder string) dateinput.Option {
	return func(opts *dateinput.Options) {
		opts.Placeholder = placeholder
	}
}

func DefaultValue(value time.Time) dateinput.Option {
	return func(opts *dateinput.Options) {
		opts.DefaultValue = &value
	}
}

func Required(required bool) dateinput.Option {
	return func(opts *dateinput.Options) {
		opts.Required = required
	}
}

func Disabled(disabled bool) dateinput.Option {
	return func(opts *dateinput.Options) {
		opts.Disabled = disabled
	}
}

func Format(format string) dateinput.Option {
	return func(opts *dateinput.Options) {
		opts.Format = format
	}
}

func MaxValue(value time.Time) dateinput.Option {
	return func(opts *dateinput.Options) {
		opts.MaxValue = &value
	}
}

func MinLength(value time.Time) dateinput.Option {
	return func(opts *dateinput.Options) {
		opts.MinValue = &value
	}
}

func Location(location time.Location) dateinput.Option {
	return func(opts *dateinput.Options) {
		opts.Location = &location
	}
}
