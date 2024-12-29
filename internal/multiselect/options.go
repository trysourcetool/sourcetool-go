package multiselect

type Options struct {
	Label        string
	Options      []string
	Placeholder  string
	DefaultValue []string
	Required     bool
	FormatFunc   func(string, int) string
}

type Option func(*Options)
