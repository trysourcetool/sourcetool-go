package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/button"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) Button(label string, opts ...button.Option) bool {
	buttonOpts := &options.ButtonOptions{
		Label:    label,
		Disabled: false,
	}

	for _, option := range opts {
		option.Apply(buttonOpts)
	}

	sess := b.session
	if sess == nil {
		return false
	}
	page := b.page
	if page == nil {
		return false
	}
	cursor := b.cursor
	if cursor == nil {
		return false
	}
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateButtonInputID(label, path)
	buttonState := sess.State.GetButton(widgetID)
	if buttonState == nil {
		buttonState = &state.ButtonState{
			ID:    widgetID,
			Value: false,
		}
	}
	buttonState.Label = buttonOpts.Label
	buttonState.Disabled = buttonOpts.Disabled
	sess.State.Set(widgetID, buttonState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeButton.String(),
		Path:       path,
		Data:       convertStateToButtonData(buttonState),
	})

	cursor.next()

	return buttonState.Value
}

func (b *uiBuilder) generateButtonInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeButton.String()+"-"+label+"-"+path.String())
}

func convertStateToButtonData(state *state.ButtonState) *websocket.ButtonData {
	return &websocket.ButtonData{
		Value:    state.Value,
		Label:    state.Label,
		Disabled: state.Disabled,
	}
}

func convertButtonDataToState(data *websocket.ButtonData) *state.ButtonState {
	if data == nil {
		return nil
	}
	return &state.ButtonState{
		Value:    data.Value,
		Label:    data.Label,
		Disabled: data.Disabled,
	}
}
