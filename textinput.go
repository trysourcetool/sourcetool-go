package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/textinput"
)

func (b *uiBuilder) TextInput(label string, opts ...textinput.Option) string {
	textInputOpts := &options.TextInputOptions{
		Label:        label,
		Placeholder:  "",
		DefaultValue: "",
		Required:     false,
		Disabled:     false,
		MaxLength:    nil,
		MinLength:    nil,
	}

	for _, option := range opts {
		option.Apply(textInputOpts)
	}

	sess := b.session
	if sess == nil {
		return ""
	}
	page := b.page
	if page == nil {
		return ""
	}
	cursor := b.cursor
	if cursor == nil {
		return ""
	}
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateTextInputID(label, path)
	textInputState := sess.State.GetTextInput(widgetID)
	if textInputState == nil {
		textInputState = &state.TextInputState{
			ID:    widgetID,
			Value: textInputOpts.DefaultValue,
		}
	}
	textInputState.Label = textInputOpts.Label
	textInputState.Placeholder = textInputOpts.Placeholder
	textInputState.DefaultValue = textInputOpts.DefaultValue
	textInputState.Required = textInputOpts.Required
	textInputState.Disabled = textInputOpts.Disabled
	textInputState.MaxLength = textInputOpts.MaxLength
	textInputState.MinLength = textInputOpts.MinLength
	sess.State.Set(widgetID, textInputState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeTextInput.String(),
		Path:       path,
		Data:       convertStateToTextInputData(textInputState),
	})

	cursor.next()

	return textInputState.Value
}

func (b *uiBuilder) generateTextInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeTextInput.String()+"-"+label+"-"+path.String())
}

func convertStateToTextInputData(state *state.TextInputState) *websocket.TextInputData {
	if state == nil {
		return nil
	}
	return &websocket.TextInputData{
		Value:        state.Value,
		Label:        state.Label,
		Placeholder:  state.Placeholder,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
		MaxLength:    state.MaxLength,
		MinLength:    state.MinLength,
	}
}

func convertTextInputDataToState(id uuid.UUID, data *websocket.TextInputData) *state.TextInputState {
	if data == nil {
		return nil
	}
	return &state.TextInputState{
		ID:           id,
		Value:        data.Value,
		Label:        data.Label,
		Placeholder:  data.Placeholder,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
		MaxLength:    data.MaxLength,
		MinLength:    data.MinLength,
	}
}
