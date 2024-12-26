package columns

import "github.com/trysourcetool/sourcetool-go/internal/columns"

func Weight(weight ...int) columns.Option {
	return func(o *columns.Options) {
		o.Weight = weight
	}
}
