package sourcetool

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/timeinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) TimeInput(label string, options ...timeinput.Option) *time.Time {
	opts := &timeinput.Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Location:     time.Local,
	}

	for _, option := range options {
		option(opts)
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
	state := sess.State.GetTimeInput(widgetID)
	if state == nil {
		state = &timeinput.State{
			ID:    widgetID,
			Value: opts.DefaultValue,
		}
	}
	state.Label = opts.Label
	state.Placeholder = opts.Placeholder
	state.DefaultValue = opts.DefaultValue
	state.Required = opts.Required
	state.Location = opts.Location
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: timeinput.WidgetType,
		Path:       path,
		Data:       convertStateToTimeInputData(state),
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateTimeInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, timeinput.WidgetType+"-"+label+"-"+path.String())
}

func convertTimeInputDataToState(id uuid.UUID, data *websocket.TimeInputData, location *time.Location) (*timeinput.State, error) {
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

	return &timeinput.State{
		ID:           id,
		Value:        value,
		Label:        data.Label,
		DefaultValue: defaultValue,
		Placeholder:  data.Placeholder,
		Required:     data.Required,
		Location:     location,
	}, nil
}

func convertStateToTimeInputData(state *timeinput.State) *websocket.TimeInputData {
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
	}
}
