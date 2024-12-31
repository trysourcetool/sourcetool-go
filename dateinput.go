package sourcetool

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) DateInput(label string, options ...dateinput.Option) *time.Time {
	opts := &dateinput.Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		Format:       "YYYY/MM/DD",
		MaxValue:     nil,
		MinValue:     nil,
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

	widgetID := b.generateDateInputID(label, path)
	state := sess.State.GetDateInput(widgetID)
	if state == nil {
		state = &dateinput.State{
			ID:    widgetID,
			Value: opts.DefaultValue,
		}
	}
	state.Label = opts.Label
	state.Placeholder = opts.Placeholder
	state.DefaultValue = opts.DefaultValue
	state.Required = opts.Required
	state.Disabled = opts.Disabled
	state.Format = opts.Format
	state.MaxValue = opts.MaxValue
	state.MinValue = opts.MinValue
	state.Location = opts.Location
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: dateinput.WidgetType,
		Path:       path,
		Data:       convertStateToDateInputData(state),
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateDateInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, dateinput.WidgetType+"-"+label+"-"+path.String())
}

func convertDateInputDataToState(id uuid.UUID, data *websocket.DateInputData, location *time.Location) (*dateinput.State, error) {
	if data == nil {
		return nil, nil
	}

	parseDate := func(dateStr string) (*time.Time, error) {
		if dateStr == "" {
			return nil, nil
		}
		t, err := time.ParseInLocation(time.DateOnly, dateStr, location)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date %q: %v", dateStr, err)
		}
		return &t, nil
	}

	value, err := parseDate(data.Value)
	if err != nil {
		return nil, err
	}

	defaultValue, err := parseDate(data.DefaultValue)
	if err != nil {
		return nil, err
	}

	maxValue, err := parseDate(data.MaxValue)
	if err != nil {
		return nil, err
	}

	minValue, err := parseDate(data.MinValue)
	if err != nil {
		return nil, err
	}

	return &dateinput.State{
		ID:           id,
		Value:        value,
		Label:        data.Label,
		DefaultValue: defaultValue,
		Placeholder:  data.Placeholder,
		Required:     data.Required,
		Disabled:     data.Disabled,
		Format:       data.Format,
		MaxValue:     maxValue,
		MinValue:     minValue,
		Location:     location,
	}, nil
}

func convertStateToDateInputData(state *dateinput.State) *websocket.DateInputData {
	if state == nil {
		return nil
	}
	var value, defaultValue, maxValue, minValue string
	if state.Value != nil {
		value = state.Value.Format(time.DateOnly)
	}
	if state.DefaultValue != nil {
		defaultValue = state.DefaultValue.Format(time.DateOnly)
	}
	if state.MaxValue != nil {
		maxValue = state.MaxValue.Format(time.DateOnly)
	}
	if state.MinValue != nil {
		minValue = state.MinValue.Format(time.DateOnly)
	}
	return &websocket.DateInputData{
		Value:        value,
		Label:        state.Label,
		Placeholder:  state.Placeholder,
		DefaultValue: defaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
		Format:       state.Format,
		MaxValue:     maxValue,
		MinValue:     minValue,
	}
}
