package table

import "github.com/gofrs/uuid/v5"

type ReturnValue struct {
	Selection *Selection `json:"selection"`
}

type Selection struct {
	Row  int   `json:"row"`
	Rows []int `json:"rows"`
}

type State struct {
	ID          uuid.UUID   `json:"-"`
	Value       ReturnValue `json:"value"`
	Header      string      `json:"header"`
	Description string      `json:"description"`
	OnSelect    string      `json:"onSelect"`
	Data        any         `json:"data"`
}
