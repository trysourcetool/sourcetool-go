package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/radio"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) Radio(label string, options ...radio.Option) *radio.Value {
	opts := &radio.Options{
		Label:        label,
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		FormatFunc:   nil,
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

	var defaultVal *int
	if opts.DefaultValue != nil {
		for i, o := range opts.Options {
			if conv.SafeValue(opts.DefaultValue) == o {
				defaultVal = &i
				break
			}
		}
	}

	widgetID := b.generateRadioID(label, path)
	state := sess.State.GetRadio(widgetID)
	if state == nil {
		state = &radio.State{
			ID:           widgetID,
			Value:        defaultVal,
			DefaultValue: defaultVal,
		}
	}

	if opts.FormatFunc == nil {
		opts.FormatFunc = func(v string, i int) string {
			return v
		}
	}

	displayVals := make([]string, len(opts.Options))
	for i, v := range opts.Options {
		displayVals[i] = opts.FormatFunc(v, i)
	}

	state.Label = opts.Label
	state.Options = displayVals
	state.DefaultValue = defaultVal
	state.Required = opts.Required
	state.Disabled = opts.Disabled
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: radio.WidgetType,
		Path:       path,
		Data:       convertStateToRadioData(state),
	})

	cursor.next()

	var value *radio.Value
	if state.Value != nil {
		value = &radio.Value{
			Value: opts.Options[*state.Value],
			Index: conv.SafeValue(state.Value),
		}
	}

	return value
}

func (b *uiBuilder) generateRadioID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, radio.WidgetType+"-"+label+"-"+path.String())
}

func convertStateToRadioData(state *radio.State) *websocket.RadioData {
	if state == nil {
		return nil
	}
	return &websocket.RadioData{
		Label:        state.Label,
		Value:        state.Value,
		Options:      state.Options,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}

func convertRadioDataToState(id uuid.UUID, data *websocket.RadioData) *radio.State {
	if data == nil {
		return nil
	}
	return &radio.State{
		ID:           id,
		Label:        data.Label,
		Value:        data.Value,
		Options:      data.Options,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
