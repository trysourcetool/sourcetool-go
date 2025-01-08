package table

import "github.com/trysourcetool/sourcetool-go/internal/options"

type Option interface {
	Apply(*options.TableOptions)
}

type headerOption string

func (h headerOption) Apply(opts *options.TableOptions) {
	opts.Header = (*string)(&h)
}

func Header(header string) Option {
	return headerOption(header)
}

type descriptionOption string

func (d descriptionOption) Apply(opts *options.TableOptions) {
	opts.Description = (*string)(&d)
}

func Description(description string) Option {
	return descriptionOption(description)
}

type onSelectOption SelectionBehavior

func (o onSelectOption) Apply(opts *options.TableOptions) {
	opts.OnSelect = (*string)((*SelectionBehavior)(&o))
}

func OnSelect(behavior SelectionBehavior) Option {
	return onSelectOption(behavior)
}

type rowSelectionOption SelectionMode

func (r rowSelectionOption) Apply(opts *options.TableOptions) {
	opts.RowSelection = (*string)((*SelectionMode)(&r))
}

func RowSelection(mode SelectionMode) Option {
	return rowSelectionOption(mode)
}
