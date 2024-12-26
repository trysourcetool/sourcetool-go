package columns

type Options struct {
	Columns int
}

type Option func(*Options)
