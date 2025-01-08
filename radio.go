package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/radio"
)

func (b *uiBuilder) Radio(label string, opts ...radio.Option) *radio.Value {
	radioOpts := &options.RadioOptions{
		Label:        label,
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		FormatFunc:   nil,
	}

	for _, option := range opts {
		option.Apply(radioOpts)
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
	if radioOpts.DefaultValue != nil {
		for i, o := range radioOpts.Options {
			if conv.SafeValue(radioOpts.DefaultValue) == o {
				defaultVal = &i
				break
			}
		}
	}

	widgetID := b.generateRadioID(label, path)
	radioState := sess.State.GetRadio(widgetID)
	if radioState == nil {
		radioState = &state.RadioState{
			ID:           widgetID,
			Value:        defaultVal,
			DefaultValue: defaultVal,
		}
	}

	if radioOpts.FormatFunc == nil {
		radioOpts.FormatFunc = func(v string, i int) string {
			return v
		}
	}

	displayVals := make([]string, len(radioOpts.Options))
	for i, v := range radioOpts.Options {
		displayVals[i] = radioOpts.FormatFunc(v, i)
	}

	radioState.Label = radioOpts.Label
	radioState.Options = displayVals
	radioState.DefaultValue = defaultVal
	radioState.Required = radioOpts.Required
	radioState.Disabled = radioOpts.Disabled
	sess.State.Set(widgetID, radioState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeRadio.String(),
		Path:       path,
		Data:       convertStateToRadioData(radioState),
	})

	cursor.next()

	var value *radio.Value
	if radioState.Value != nil {
		value = &radio.Value{
			Value: radioOpts.Options[*radioState.Value],
			Index: conv.SafeValue(radioState.Value),
		}
	}

	return value
}

func (b *uiBuilder) generateRadioID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeRadio.String()+"-"+label+"-"+path.String())
}

func convertStateToRadioData(state *state.RadioState) *websocket.RadioData {
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

func convertRadioDataToState(id uuid.UUID, data *websocket.RadioData) *state.RadioState {
	if data == nil {
		return nil
	}
	return &state.RadioState{
		ID:           id,
		Label:        data.Label,
		Value:        data.Value,
		Options:      data.Options,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
