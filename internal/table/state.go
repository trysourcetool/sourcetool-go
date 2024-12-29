package table

import "github.com/gofrs/uuid/v5"

const WidgetType = "table"

type Value struct {
	Selection *Selection
}

type Selection struct {
	Row  int
	Rows []int
}

type State struct {
	ID           uuid.UUID
	Data         any
	Value        Value
	Header       string
	Description  string
	OnSelect     string
	RowSelection string
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
