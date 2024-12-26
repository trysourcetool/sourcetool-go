package columns

type Options struct {
	Columns int
	Weight  []int
}

type Option func(*Options)
