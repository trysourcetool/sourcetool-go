package session

import (
	"sync"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/button"
	"github.com/trysourcetool/sourcetool-go/internal/checkbox"
	"github.com/trysourcetool/sourcetool-go/internal/checkboxgroup"
	"github.com/trysourcetool/sourcetool-go/internal/columns"
	"github.com/trysourcetool/sourcetool-go/internal/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/datetimeinput"
	"github.com/trysourcetool/sourcetool-go/internal/form"
	"github.com/trysourcetool/sourcetool-go/internal/markdown"
	"github.com/trysourcetool/sourcetool-go/internal/multiselect"
	"github.com/trysourcetool/sourcetool-go/internal/numberinput"
	"github.com/trysourcetool/sourcetool-go/internal/radio"
	"github.com/trysourcetool/sourcetool-go/internal/selectbox"
	"github.com/trysourcetool/sourcetool-go/internal/table"
	"github.com/trysourcetool/sourcetool-go/internal/textarea"
	"github.com/trysourcetool/sourcetool-go/internal/textinput"
	"github.com/trysourcetool/sourcetool-go/internal/timeinput"
)

type WidgetState interface {
	IsWidgetState()
	GetType() string
}

type StateData map[uuid.UUID]WidgetState

type State struct {
	// data map[uuid.UUID]any // ui ID -> options state
	data StateData
	mu   sync.RWMutex
}

func newState() *State {
	return &State{
		data: make(map[uuid.UUID]WidgetState),
	}
}

func (s *State) Get(id uuid.UUID) WidgetState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return state
}

func (s *State) GetTextInput(id uuid.UUID) *textinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*textinput.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetNumberInput(id uuid.UUID) *numberinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*numberinput.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetDateInput(id uuid.UUID) *dateinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*dateinput.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetDateTimeInput(id uuid.UUID) *datetimeinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*datetimeinput.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetTimeInput(id uuid.UUID) *timeinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*timeinput.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetSelectbox(id uuid.UUID) *selectbox.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*selectbox.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetMultiSelect(id uuid.UUID) *multiselect.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*multiselect.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetCheckbox(id uuid.UUID) *checkbox.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*checkbox.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetCheckboxGroup(id uuid.UUID) *checkboxgroup.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*checkboxgroup.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetRadio(id uuid.UUID) *radio.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*radio.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetTextArea(id uuid.UUID) *textarea.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*textarea.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetTable(id uuid.UUID) *table.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*table.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetButton(id uuid.UUID) *button.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*button.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetColumns(id uuid.UUID) *columns.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*columns.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetMarkdown(id uuid.UUID) *markdown.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*markdown.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetForm(id uuid.UUID) *form.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := state.(*form.State)
	if !ok {
		return nil
	}

	return v
}

func (s *State) ResetStates() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[uuid.UUID]WidgetState)
}

func (s *State) ResetButtons() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, state := range s.data {
		switch state.GetType() {
		case button.WidgetType:
			buttonState, ok := state.(*button.State)
			if ok {
				buttonState.Value = false
				s.data[id] = buttonState
			}
		case form.WidgetType:
			formState, ok := state.(*form.State)
			if ok {
				formState.Value = false
				s.data[id] = formState
			}
		}
	}
}

func (s *State) Set(id uuid.UUID, state WidgetState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = state
}

func (s *State) SetStates(states map[uuid.UUID]WidgetState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, state := range states {
		s.data[id] = state
	}
}
