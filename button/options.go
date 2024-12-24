package button

import "github.com/trysourcetool/sourcetool-go/internal/button"

func Disabled(disabled bool) button.Option {
	return func(opts *button.Options) {
		opts.Disabled = disabled
	}
}
