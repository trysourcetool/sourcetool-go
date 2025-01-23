package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/numberinput"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
)

func (b *uiBuilder) NumberInput(label string, opts ...numberinput.Option) *float64 {
	numberInputOpts := &options.NumberInputOptions{
		Label:        label,
		Placeholder:  "",
		DefaultValue: conv.NilValue(float64(0)),
		Required:     false,
		Disabled:     false,
		MaxValue:     nil,
		MinValue:     nil,
	}

	for _, o := range opts {
		o.Apply(numberInputOpts)
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

	widgetID := b.generateNumberInputID(label, path)
	numberInputState := sess.State.GetNumberInput(widgetID)
	if numberInputState == nil {
		numberInputState = &state.NumberInputState{
			ID:    widgetID,
			Value: numberInputOpts.DefaultValue,
		}
	}
	numberInputState.Label = numberInputOpts.Label
	numberInputState.Placeholder = numberInputOpts.Placeholder
	numberInputState.DefaultValue = numberInputOpts.DefaultValue
	numberInputState.Required = numberInputOpts.Required
	numberInputState.Disabled = numberInputOpts.Disabled
	numberInputState.MaxValue = numberInputOpts.MaxValue
	numberInputState.MinValue = numberInputOpts.MinValue
	sess.State.Set(widgetID, numberInputState)

	numberInput := convertStateToNumberInputProto(numberInputState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_NumberInput{
				NumberInput: numberInput,
			},
		},
	})

	cursor.next()

	return numberInputState.Value
}

func (b *uiBuilder) generateNumberInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeNumberInput.String()+"-"+label+"-"+path.String())
}

func convertStateToNumberInputProto(state *state.NumberInputState) *widgetv1.NumberInput {
	if state == nil {
		return nil
	}
	return &widgetv1.NumberInput{
		Value:        state.Value,
		Label:        state.Label,
		Placeholder:  state.Placeholder,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
		MaxValue:     state.MaxValue,
		MinValue:     state.MinValue,
	}
}

func convertNumberInputProtoToState(id uuid.UUID, data *widgetv1.NumberInput) *state.NumberInputState {
	if data == nil {
		return nil
	}
	return &state.NumberInputState{
		ID:           id,
		Value:        data.Value,
		Label:        data.Label,
		Placeholder:  data.Placeholder,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
		MaxValue:     data.MaxValue,
		MinValue:     data.MinValue,
	}
}
