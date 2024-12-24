package table

type OnSelect string

func (o OnSelect) String() string {
	return string(o)
}

type RowSelection string

func (r RowSelection) String() string {
	return string(r)
}

type Options struct {
	Header       string
	Description  string
	OnSelect     OnSelect
	RowSelection RowSelection
}

type Option func(*Options)
