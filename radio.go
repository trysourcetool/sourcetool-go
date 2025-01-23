package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/radio"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
)

func (b *uiBuilder) Radio(label string, opts ...radio.Option) *radio.Value {
	radioOpts := &options.RadioOptions{
		Label:        label,
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		FormatFunc:   nil,
	}

	for _, o := range opts {
		o.Apply(radioOpts)
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

	radioProto := convertStateToRadioProto(radioState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_Radio{
				Radio: radioProto,
			},
		},
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

func convertStateToRadioProto(state *state.RadioState) *widgetv1.Radio {
	if state == nil {
		return nil
	}
	var value *int64
	if state.Value != nil {
		val := int64(*state.Value)
		value = &val
	}
	var defaultValue *int64
	if state.DefaultValue != nil {
		val := int64(*state.DefaultValue)
		defaultValue = &val
	}
	return &widgetv1.Radio{
		Label:        state.Label,
		Value:        value,
		Options:      state.Options,
		DefaultValue: defaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}

func convertRadioProtoToState(id uuid.UUID, data *widgetv1.Radio) *state.RadioState {
	if data == nil {
		return nil
	}
	var value *int
	if data.Value != nil {
		val := int(*data.Value)
		value = &val
	}
	var defaultValue *int
	if data.DefaultValue != nil {
		val := int(*data.DefaultValue)
		defaultValue = &val
	}
	return &state.RadioState{
		ID:           id,
		Label:        data.Label,
		Value:        value,
		Options:      data.Options,
		DefaultValue: defaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
