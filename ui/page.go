package ui

type Page struct {
	ID      string
	Name    string
	Handler func() error
}
