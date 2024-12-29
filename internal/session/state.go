package session

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/button"
	"github.com/trysourcetool/sourcetool-go/internal/columns"
	"github.com/trysourcetool/sourcetool-go/internal/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/form"
	"github.com/trysourcetool/sourcetool-go/internal/markdown"
	"github.com/trysourcetool/sourcetool-go/internal/multiselect"
	"github.com/trysourcetool/sourcetool-go/internal/numberinput"
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

func (s *State) GetDateInput(id uuid.UUID) *dateinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToDateInputState(state)
}

func anyToDateInputState(a any) *dateinput.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var dateInputState dateinput.State
	if err := json.Unmarshal(bytes, &dateInputState); err != nil {
		return nil
	}

	return &dateInputState
}

func (s *State) GetTimeInput(id uuid.UUID) *timeinput.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToTimeInputState(state)
}

func anyToTimeInputState(a any) *timeinput.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var timeInputState timeinput.State
	if err := json.Unmarshal(bytes, &timeInputState); err != nil {
		return nil
	}

	return &timeInputState
}

func (s *State) GetSelectbox(id uuid.UUID) *selectbox.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToSelectboxState(state)
}

func anyToSelectboxState(a any) *selectbox.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var selectboxState selectbox.State
	if err := json.Unmarshal(bytes, &selectboxState); err != nil {
		return nil
	}

	return &selectboxState
}

func (s *State) GetMultiSelect(id uuid.UUID) *multiselect.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return anyToMultiSelectState(state)
}

func anyToMultiSelectState(a any) *multiselect.State {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil
	}

	var multiSelectState multiselect.State
	if err := json.Unmarshal(bytes, &multiSelectState); err != nil {
		return nil
	}

	return &multiSelectState
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
	s.data = make(map[uuid.UUID]WidgetState)
}

func (s *State) ResetButtons() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, state := range s.data {
		switch state.GetType() {
		case button.WidgetType:
			buttonState := anyToButtonState(state)
			buttonState.Value = false
		case form.WidgetType:
			formState := anyToFormState(state)
			formState.Value = false
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

type widgetStateJSON struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (s *StateData) UnmarshalJSON(data []byte) error {
	var raw map[uuid.UUID]widgetStateJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*s = make(StateData)
	for id, stateData := range raw {
		var state WidgetState
		switch stateData.Type {
		case "textinput":
			state = new(textinput.State)
		case "numberinput":
			state = new(numberinput.State)
		case "dateinput":
			state = new(dateinput.State)
		// 他の型も同様に追加
		default:
			return fmt.Errorf("unknown widget type: %s", stateData.Type)
		}

		if err := json.Unmarshal(stateData.Value, state); err != nil {
			return err
		}
		(*s)[id] = state
	}
	return nil
}

// MarshalJSONも実装する必要があります
func (s StateData) MarshalJSON() ([]byte, error) {
	result := make(map[uuid.UUID]widgetStateJSON)
	for id, state := range s {
		var stateType string
		switch state.(type) {
		case *textinput.State:
			stateType = "textinput"
		case *numberinput.State:
			stateType = "numberinput"
		case *dateinput.State:
			stateType = "dateinput"
		// 他の型も同様に追加
		default:
			return nil, fmt.Errorf("unknown widget type")
		}

		value, err := json.Marshal(state)
		if err != nil {
			return nil, err
		}

		result[id] = widgetStateJSON{
			Type:  stateType,
			Value: value,
		}
	}
	return json.Marshal(result)
}
