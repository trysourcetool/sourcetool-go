package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/selectbox"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) Selectbox(label string, options ...selectbox.Option) *any {
	opts := &selectbox.Options{
		Label:        label,
		DefaultIndex: nil,
		Placeholder:  "",
		Required:     false,
		DisplayFunc:  nil,
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

	if opts.DisplayFunc == nil {
		opts.DisplayFunc = func(v any, i int) string {
			return conv.StringValue(v)
		}
	}

	displayVals := make([]string, len(opts.Options))
	for i, v := range opts.Options {
		displayVals[i] = opts.DisplayFunc(v, i)
	}

	widgetID := b.generateSelectboxID(label, path)
	state := sess.State.GetSelectbox(widgetID)
	if state == nil {
		var val any
		if opts.DefaultIndex != nil && len(opts.Options) > conv.SafeValue(opts.DefaultIndex) {
			val = opts.Options[*opts.DefaultIndex]
		}
		state = &selectbox.State{
			ID:    widgetID,
			Value: conv.NilValue(val),
		}
	}
	state.Labels = displayVals
	state.Label = opts.Label
	state.Placeholder = opts.Placeholder
	state.DefaultIndex = opts.DefaultIndex
	state.Required = opts.Required
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: selectbox.WidgetType,
		Path:       path,
		Data:       convertStateToSelectboxData(state),
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateSelectboxID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, selectbox.WidgetType+"-"+label+"-"+path.String())
}

func convertStateToSelectboxData(state *selectbox.State) *websocket.SelectboxData {
	if state == nil {
		return nil
	}
	return &websocket.SelectboxData{
		Value:        state.Value,
		Values:       state.Values,
		Labels:       state.Labels,
		Label:        state.Label,
		Placeholder:  state.Placeholder,
		DefaultIndex: state.DefaultIndex,
		Required:     state.Required,
	}
}

func convertSelectboxDataToState(id uuid.UUID, data *websocket.SelectboxData) *selectbox.State {
	if data == nil {
		return nil
	}
	return &selectbox.State{
		ID:           id,
		Value:        data.Value,
		Values:       data.Values,
		Labels:       data.Labels,
		Label:        data.Label,
		Placeholder:  data.Placeholder,
		DefaultIndex: data.DefaultIndex,
		Required:     data.Required,
	}
}
