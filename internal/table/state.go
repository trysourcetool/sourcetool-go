package table

import "github.com/gofrs/uuid/v5"

const WidgetType = "table"

type Value struct {
	Selection *Selection `json:"selection"`
}

type Selection struct {
	Row  int   `json:"row"`
	Rows []int `json:"rows"`
}

type State struct {
	ID           uuid.UUID `json:"-"`
	Data         any       `json:"data"`
	Value        Value     `json:"value"`
	Header       string    `json:"header"`
	Description  string    `json:"description"`
	OnSelect     string    `json:"onSelect"`
	RowSelection string    `json:"rowSelection"`
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
