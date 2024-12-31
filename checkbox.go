package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/checkbox"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) Checkbox(label string, options ...checkbox.Option) bool {
	opts := &checkbox.Options{
		Label:        label,
		DefaultValue: false,
		Required:     false,
		Disabled:     false,
	}

	for _, option := range options {
		option(opts)
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

	widgetID := b.generateCheckboxID(label, path)
	state := sess.State.GetCheckbox(widgetID)
	if state == nil {
		state = &checkbox.State{
			ID:    widgetID,
			Value: opts.DefaultValue,
		}
	}
	state.Label = opts.Label
	state.DefaultValue = opts.DefaultValue
	state.Required = opts.Required
	state.Disabled = opts.Disabled
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: checkbox.WidgetType,
		Path:       path,
		Data:       convertStateToCheckboxData(state),
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateCheckboxID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, checkbox.WidgetType+"-"+label+"-"+path.String())
}

func convertStateToCheckboxData(state *checkbox.State) *websocket.CheckboxData {
	if state == nil {
		return nil
	}
	return &websocket.CheckboxData{
		Value:        state.Value,
		Label:        state.Label,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}

func convertCheckboxDataToState(id uuid.UUID, data *websocket.CheckboxData) *checkbox.State {
	if data == nil {
		return nil
	}
	return &checkbox.State{
		ID:           id,
		Value:        data.Value,
		Label:        data.Label,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
