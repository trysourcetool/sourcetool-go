package radio

type Options struct {
	Label        string
	Options      []string
	DefaultValue *string
	Required     bool
	Disabled     bool
	FormatFunc   func(string, int) string
}

type Option func(*Options)
