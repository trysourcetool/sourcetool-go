package table

type onSelect string

const (
	OnSelectIgnore onSelect = "ignore"
	OnSelectRerun  onSelect = "rerun"
)

func (o onSelect) String() string {
	return string(o)
}

type Options struct {
	Header      string
	Description string
	OnSelect    onSelect
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
