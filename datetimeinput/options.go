package datetimeinput

import (
	"time"

	"github.com/trysourcetool/sourcetool-go/internal/datetimeinput"
)

func Placeholder(placeholder string) datetimeinput.Option {
	return func(opts *datetimeinput.Options) {
		opts.Placeholder = placeholder
	}
}

func DefaultValue(value time.Time) datetimeinput.Option {
	return func(opts *datetimeinput.Options) {
		opts.DefaultValue = &value
	}
}

func Required(required bool) datetimeinput.Option {
	return func(opts *datetimeinput.Options) {
		opts.Required = required
	}
}

func Disabled(disabled bool) datetimeinput.Option {
	return func(opts *datetimeinput.Options) {
		opts.Disabled = disabled
	}
}

func Format(format string) datetimeinput.Option {
	return func(opts *datetimeinput.Options) {
		opts.Format = format
	}
}

func MaxValue(value time.Time) datetimeinput.Option {
	return func(opts *datetimeinput.Options) {
		opts.MaxValue = &value
	}
}

func MinLength(value time.Time) datetimeinput.Option {
	return func(opts *datetimeinput.Options) {
		opts.MinValue = &value
	}
}

func Location(location time.Location) datetimeinput.Option {
	return func(opts *datetimeinput.Options) {
		opts.Location = &location
	}
}
