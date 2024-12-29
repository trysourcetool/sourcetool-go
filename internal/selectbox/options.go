package selectbox

type Options struct {
	Label        string
	Options      []any
	Placeholder  string
	DefaultIndex *int
	Required     bool
	DisplayFunc  func(any, int) string
}

type Option func(*Options)
