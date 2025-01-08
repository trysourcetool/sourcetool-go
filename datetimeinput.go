package sourcetool

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/datetimeinput"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) DateTimeInput(label string, opts ...datetimeinput.Option) *time.Time {
	dateTimeInputOpts := &options.DateTimeInputOptions{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		Format:       "YYYY/MM/DD HH:MM:SS",
		MaxValue:     nil,
		MinValue:     nil,
		Location:     time.Local,
	}

	for _, o := range opts {
		o.Apply(dateTimeInputOpts)
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
	dateTimeInputState := sess.State.GetDateTimeInput(widgetID)
	if dateTimeInputState == nil {
		dateTimeInputState = &state.DateTimeInputState{
			ID:    widgetID,
			Value: dateTimeInputOpts.DefaultValue,
		}
	}
	dateTimeInputState.Label = dateTimeInputOpts.Label
	dateTimeInputState.Placeholder = dateTimeInputOpts.Placeholder
	dateTimeInputState.DefaultValue = dateTimeInputOpts.DefaultValue
	dateTimeInputState.Required = dateTimeInputOpts.Required
	dateTimeInputState.Disabled = dateTimeInputOpts.Disabled
	dateTimeInputState.Format = dateTimeInputOpts.Format
	dateTimeInputState.MaxValue = dateTimeInputOpts.MaxValue
	dateTimeInputState.MinValue = dateTimeInputOpts.MinValue
	dateTimeInputState.Location = dateTimeInputOpts.Location
	sess.State.Set(widgetID, dateTimeInputState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeDateTimeInput.String(),
		Path:       path,
		Data:       convertStateToDateTimeInputData(dateTimeInputState),
	})

	cursor.next()

	return dateTimeInputState.Value
}

func (b *uiBuilder) generateDateTimeInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeDateTimeInput.String()+"-"+label+"-"+path.String())
}

func convertDateTimeInputDataToState(id uuid.UUID, data *websocket.DateTimeInputData, location *time.Location) (*state.DateTimeInputState, error) {
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

	return &state.DateTimeInputState{
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

func convertStateToDateTimeInputData(state *state.DateTimeInputState) *websocket.DateTimeInputData {
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
		Disabled:     state.Disabled,
		Format:       state.Format,
		MaxValue:     maxValue,
		MinValue:     minValue,
	}
}
