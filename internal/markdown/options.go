package markdown

type Options struct {
	Body string
}

type Option func(*Options)
