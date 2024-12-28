package timeinput

import (
	"time"

	"github.com/trysourcetool/sourcetool-go/internal/timeinput"
)

func Placeholder(placeholder string) timeinput.Option {
	return func(opts *timeinput.Options) {
		opts.Placeholder = placeholder
	}
}

func DefaultValue(value time.Time) timeinput.Option {
	return func(opts *timeinput.Options) {
		opts.DefaultValue = &value
	}
}

func Required(required bool) timeinput.Option {
	return func(opts *timeinput.Options) {
		opts.Required = required
	}
}

func Location(location time.Location) timeinput.Option {
	return func(opts *timeinput.Options) {
		opts.Location = &location
	}
}
