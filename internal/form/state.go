package form

import "github.com/gofrs/uuid/v5"

type State struct {
	ID             uuid.UUID `json:"-"`
	Value          bool      `json:"value"`
	ButtonLabel    string    `json:"buttonLabel"`
	ButtonDisabled bool      `json:"buttonDisabled"`
	ClearOnSubmit  bool      `json:"clearOnSubmit"`
}
