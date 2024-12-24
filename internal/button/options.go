package button

type Options struct {
	Label    string
	Disabled bool
}

type Option func(*Options)
