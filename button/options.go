package button

type Options struct {
	Label    string
	Disabled bool
}

type Option func(*Options)

func Disabled(disabled bool) Option {
	return func(opts *Options) {
		opts.Disabled = disabled
	}
}
