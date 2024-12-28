package markdown

import "github.com/gofrs/uuid/v5"

const WidgetType = "markdown"

type State struct {
	ID   uuid.UUID
	Body string
}

func (s *State) IsWidgetState()  {}
func (s *State) GetType() string { return WidgetType }
