package sourcetool

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/timeinput"
)

func (b *uiBuilder) TimeInput(label string, opts ...timeinput.Option) *time.Time {
	timeInputOpts := &options.TimeInputOptions{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		Location:     time.Local,
	}

	for _, option := range opts {
		option.Apply(timeInputOpts)
	}

	sess := b.session
	if sess == nil {
		return nil
	}
	page := b.page
	if page == nil {
		return nil
	}
	cursor := b.cursor
	if cursor == nil {
		return nil
	}
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateTimeInputID(label, path)
	timeInputState := sess.State.GetTimeInput(widgetID)
	if timeInputState == nil {
		timeInputState = &state.TimeInputState{
			ID:    widgetID,
			Value: timeInputOpts.DefaultValue,
		}
	}
	timeInputState.Label = timeInputOpts.Label
	timeInputState.Placeholder = timeInputOpts.Placeholder
	timeInputState.DefaultValue = timeInputOpts.DefaultValue
	timeInputState.Required = timeInputOpts.Required
	timeInputState.Disabled = timeInputOpts.Disabled
	timeInputState.Location = timeInputOpts.Location
	sess.State.Set(widgetID, timeInputState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeTimeInput.String(),
		Path:       path,
		Data:       convertStateToTimeInputData(timeInputState),
	})

	cursor.next()

	return timeInputState.Value
}

func (b *uiBuilder) generateTimeInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeTimeInput.String()+"-"+label+"-"+path.String())
}

func convertTimeInputDataToState(id uuid.UUID, data *websocket.TimeInputData, location *time.Location) (*state.TimeInputState, error) {
	if data == nil {
		return nil, nil
	}

	parseTime := func(timeStr string) (*time.Time, error) {
		if timeStr == "" {
			return nil, nil
		}
		t, err := time.ParseInLocation(time.TimeOnly, timeStr, location)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time %q: %v", timeStr, err)
		}
		return &t, nil
	}

	value, err := parseTime(data.Value)
	if err != nil {
		return nil, err
	}

	defaultValue, err := parseTime(data.DefaultValue)
	if err != nil {
		return nil, err
	}

	return &state.TimeInputState{
		ID:           id,
		Value:        value,
		Label:        data.Label,
		DefaultValue: defaultValue,
		Placeholder:  data.Placeholder,
		Required:     data.Required,
		Disabled:     data.Disabled,
		Location:     location,
	}, nil
}

func convertStateToTimeInputData(state *state.TimeInputState) *websocket.TimeInputData {
	if state == nil {
		return nil
	}
	var value, defaultValue string
	if state.Value != nil {
		value = state.Value.Format(time.TimeOnly)
	}
	if state.DefaultValue != nil {
		defaultValue = state.DefaultValue.Format(time.TimeOnly)
	}
	return &websocket.TimeInputData{
		Value:        value,
		Label:        state.Label,
		Placeholder:  state.Placeholder,
		DefaultValue: defaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}
