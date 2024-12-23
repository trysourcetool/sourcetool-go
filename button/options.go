package button

type Options struct {
	Label    string
	Disabled bool
}

func DefaultOptions(label string) *Options {
	return &Options{
		Label:    label,
		Disabled: false,
	}
}

type Option func(*Options)

func WithDisabled(disabled bool) Option {
	return func(opts *Options) {
		opts.Disabled = disabled
	}
}
