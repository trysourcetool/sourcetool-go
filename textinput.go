package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/textinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) TextInput(label string, options ...textinput.Option) string {
	opts := &textinput.Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: "",
		Required:     false,
		Disabled:     false,
		MaxLength:    nil,
		MinLength:    nil,
	}

	for _, option := range options {
		option(opts)
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
	state := sess.State.GetTextInput(widgetID)
	if state == nil {
		state = &textinput.State{
			ID:    widgetID,
			Value: opts.DefaultValue,
		}
	}
	state.Label = opts.Label
	state.Placeholder = opts.Placeholder
	state.DefaultValue = opts.DefaultValue
	state.Required = opts.Required
	state.Disabled = opts.Disabled
	state.MaxLength = opts.MaxLength
	state.MinLength = opts.MinLength
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: textinput.WidgetType,
		Path:       path,
		Data:       convertStateToTextInputData(state),
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateTextInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, textinput.WidgetType+"-"+label+"-"+path.String())
}

func convertStateToTextInputData(state *textinput.State) *websocket.TextInputData {
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

func convertTextInputDataToState(id uuid.UUID, data *websocket.TextInputData) *textinput.State {
	if data == nil {
		return nil
	}
	return &textinput.State{
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
