package sourcetool

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/datetimeinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) DateTimeInput(label string, options ...datetimeinput.Option) *time.Time {
	opts := &datetimeinput.Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Format:       "YYYY/MM/DD HH:MM:SS",
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

	widgetID := b.generateDateTimeInputID(label, path)
	state := sess.State.GetDateTimeInput(widgetID)
	if state == nil {
		state = &datetimeinput.State{
			ID:    widgetID,
			Value: opts.DefaultValue,
		}
	}
	state.Label = opts.Label
	state.Placeholder = opts.Placeholder
	state.DefaultValue = opts.DefaultValue
	state.Required = opts.Required
	state.Format = opts.Format
	state.MaxValue = opts.MaxValue
	state.MinValue = opts.MinValue
	state.Location = opts.Location
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: datetimeinput.WidgetType,
		Path:       path,
		Data:       convertStateToDateTimeInputData(state),
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateDateTimeInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, datetimeinput.WidgetType+"-"+label+"-"+path.String())
}

func convertDateTimeInputDataToState(id uuid.UUID, data *websocket.DateTimeInputData, location *time.Location) (*datetimeinput.State, error) {
	if data == nil {
		return nil, nil
	}

	parseDate := func(dateStr string) (*time.Time, error) {
		if dateStr == "" {
			return nil, nil
		}
		t, err := time.ParseInLocation(time.DateTime, dateStr, location)
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

	return &datetimeinput.State{
		ID:           id,
		Value:        value,
		Label:        data.Label,
		DefaultValue: defaultValue,
		Placeholder:  data.Placeholder,
		Required:     data.Required,
		Format:       data.Format,
		MaxValue:     maxValue,
		MinValue:     minValue,
		Location:     location,
	}, nil
}

func convertStateToDateTimeInputData(state *datetimeinput.State) *websocket.DateTimeInputData {
	if state == nil {
		return nil
	}
	var value, defaultValue, maxValue, minValue string
	if state.Value != nil {
		value = state.Value.Format(time.DateTime)
	}
	if state.DefaultValue != nil {
		defaultValue = state.DefaultValue.Format(time.DateTime)
	}
	if state.MaxValue != nil {
		maxValue = state.MaxValue.Format(time.DateTime)
	}
	if state.MinValue != nil {
		minValue = state.MinValue.Format(time.DateTime)
	}
	return &websocket.DateTimeInputData{
		Value:        value,
		Label:        state.Label,
		Placeholder:  state.Placeholder,
		DefaultValue: defaultValue,
		Required:     state.Required,
		Format:       state.Format,
		MaxValue:     maxValue,
		MinValue:     minValue,
	}
}
