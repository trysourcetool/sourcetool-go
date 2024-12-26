package form

import "github.com/gofrs/uuid/v5"

type ReturnValue bool

type State struct {
	ID             uuid.UUID   `json:"-"`
	Value          ReturnValue `json:"value"`
	ButtonLabel    string      `json:"buttonLabel"`
	ButtonDisabled bool        `json:"buttonDisabled"`
	ClearOnSubmit  bool        `json:"clearOnSubmit"`
}
