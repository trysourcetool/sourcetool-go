package form

import "github.com/trysourcetool/sourcetool-go/internal/form"

func ClearOnSubmit(clearOnSubmit bool) form.Option {
	return func(o *form.Options) {
		o.ClearOnSubmit = clearOnSubmit
	}
}
