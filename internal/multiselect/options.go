package multiselect

type Options struct {
	Label        string
	Options      []string
	Placeholder  string
	DefaultValue []string
	Required     bool
	DisplayFunc  func(string, int) string
}

type Option func(*Options)
