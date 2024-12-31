package checkbox

import "github.com/trysourcetool/sourcetool-go/internal/checkbox"

func DefaultValue(defaultValue bool) checkbox.Option {
	return func(o *checkbox.Options) {
		o.DefaultValue = defaultValue
	}
}

func Required(required bool) checkbox.Option {
	return func(o *checkbox.Options) {
		o.Required = required
	}
}

func Disabled(disabled bool) checkbox.Option {
	return func(o *checkbox.Options) {
		o.Disabled = disabled
	}
}
