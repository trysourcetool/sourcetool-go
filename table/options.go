package table

import "github.com/trysourcetool/sourcetool-go/internal/table"

const (
	OnSelectIgnore table.OnSelect = "ignore"
	OnSelectRerun  table.OnSelect = "rerun"
)

const (
	RowSelectionSingle   table.RowSelection = "single"
	RowSelectionMultiple table.RowSelection = "multiple"
)

func Header(header string) table.Option {
	return func(opts *table.Options) {
		opts.Header = header
	}
}

func Description(description string) table.Option {
	return func(opts *table.Options) {
		opts.Description = description
	}
}

func OnSelect(onSelect table.OnSelect) table.Option {
	return func(opts *table.Options) {
		opts.OnSelect = onSelect
	}
}

func RowSelection(rowSelection table.RowSelection) table.Option {
	return func(opts *table.Options) {
		opts.RowSelection = rowSelection
	}
}
