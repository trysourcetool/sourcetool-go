package session

import (
	"encoding/json"
	"sync"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/button"
	"github.com/trysourcetool/sourcetool-go/internal/columns"
	"github.com/trysourcetool/sourcetool-go/internal/form"
	"github.com/trysourcetool/sourcetool-go/internal/markdown"
	"github.com/trysourcetool/sourcetool-go/internal/numberinput"
	"github.com/trysourcetool/sourcetool-go/internal/table"
	"github.com/trysourcetool/sourcetool-go/internal/textarea"
	"github.com/trysourcetool/sourcetool-go/internal/textinput"
)

type State struct {
	data map[uuid.UUID]any // ui ID -> options state
	mu   sync.RWMutex
}

func newState() *State {
	return &State{
		data: make(map[uuid.UUID]any),
	}
}

func (s *State) GetTextInput(id uuid.UUID) *textinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToTextInputState(state)
}

func anyToTextInputState(a any) *textinput.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var textInputState textinput.State
	if err := json.Unmarshal(bytes, &textInputState); err != nil {
		return nil
	}

	return &textInputState
}

func (s *State) GetNumberInput(id uuid.UUID) *numberinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToNumberInputState(state)
}

func anyToNumberInputState(a any) *numberinput.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var numberInputState numberinput.State
	if err := json.Unmarshal(bytes, &numberInputState); err != nil {
		return nil
	}

	return &numberInputState
}

func (s *State) GetTextArea(id uuid.UUID) *textarea.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToTextAreaState(state)
}

func anyToTextAreaState(a any) *textarea.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var textAreaState textarea.State
	if err := json.Unmarshal(bytes, &textAreaState); err != nil {
		return nil
	}

	return &textAreaState
}

func (s *State) GetTable(id uuid.UUID) *table.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToTableState(state)
}

func anyToTableState(a any) *table.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var tableState table.State
	if err := json.Unmarshal(bytes, &tableState); err != nil {
		return nil
	}

	return &tableState
}

func (s *State) GetButton(id uuid.UUID) *button.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToButtonState(state)
}

func anyToButtonState(a any) *button.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var buttonState button.State
	if err := json.Unmarshal(bytes, &buttonState); err != nil {
		return nil
	}

	return &buttonState
}

func (s *State) GetColumns(id uuid.UUID) *columns.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToColumnsState(state)
}

func anyToColumnsState(a any) *columns.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var columnsState columns.State
	if err := json.Unmarshal(bytes, &columnsState); err != nil {
		return nil
	}

	return &columnsState
}

func (s *State) GetMarkdown(id uuid.UUID) *markdown.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToMarkdownState(state)
}

func anyToMarkdownState(a any) *markdown.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var markdownState markdown.State
	if err := json.Unmarshal(bytes, &markdownState); err != nil {
		return nil
	}

	return &markdownState
}

func (s *State) GetForm(id uuid.UUID) *form.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToFormState(state)
}

func anyToFormState(a any) *form.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var formState form.State
	if err := json.Unmarshal(bytes, &formState); err != nil {
		return nil
	}

	return &formState
}

func (s *State) ResetStates() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[uuid.UUID]any)
}

func (s *State) ResetButtons() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, state := range s.data {
		buttonState := anyToButtonState(state)
		if buttonState != nil {
			buttonState.Value = false
		}
	}
}

func (s *State) Set(id uuid.UUID, state any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = state
}

func (s *State) SetStates(states map[uuid.UUID]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, state := range states {
		s.data[id] = state
	}
}
