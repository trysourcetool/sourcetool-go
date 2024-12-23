package table

type onSelect string

const (
	OnSelectIgnore onSelect = "ignore"
	OnSelectRerun  onSelect = "rerun"
)

func (o onSelect) String() string {
	return string(o)
}

type rowSelection string

const (
	RowSelectionSingle   rowSelection = "single"
	RowSelectionMultiple rowSelection = "multiple"
)

func (r rowSelection) String() string {
	return string(r)
}

type Options struct {
	Header       string
	Description  string
	OnSelect     onSelect
	RowSelection rowSelection
}

func DefaultOptions() *Options {
	return &Options{
		Header:       "",
		Description:  "",
		OnSelect:     OnSelectIgnore,
		RowSelection: RowSelectionSingle,
	}
}

type Option func(*Options)

func WithHeader(header string) Option {
	return func(opts *Options) {
		opts.Header = header
	}
}

func WithDescription(description string) Option {
	return func(opts *Options) {
		opts.Description = description
	}
}

func WithOnSelect(onSelect onSelect) Option {
	return func(opts *Options) {
		opts.OnSelect = onSelect
	}
}

func WithRowSelection(rowSelection rowSelection) Option {
	return func(opts *Options) {
		opts.RowSelection = rowSelection
	}
}
