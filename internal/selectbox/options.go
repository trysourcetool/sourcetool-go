package selectbox

type Options struct {
	Label        string
	Options      []string
	Placeholder  string
	DefaultIndex *int
	Required     bool
	DisplayFunc  func(string, int) string
}

type Option func(*Options)
