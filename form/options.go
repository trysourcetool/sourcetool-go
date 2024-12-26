package form

import "github.com/trysourcetool/sourcetool-go/internal/form"

func ButtonDisabled(buttonDisabled bool) form.Option {
	return func(o *form.Options) {
		o.ButtonDisabled = buttonDisabled
	}
}

func ClearOnSubmit(clearOnSubmit bool) form.Option {
	return func(o *form.Options) {
		o.ClearOnSubmit = clearOnSubmit
	}
}
