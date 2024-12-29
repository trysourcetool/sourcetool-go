package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/selectbox"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) Selectbox(label string, options ...selectbox.Option) *selectbox.Value {
	opts := &selectbox.Options{
		Label:        label,
		DefaultValue: nil,
		Placeholder:  "",
		Required:     false,
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

	widgetID := b.generateSelectboxID(label, path)
	state := sess.State.GetSelectbox(widgetID)
	if state == nil {
		state = &selectbox.State{
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
	state.Placeholder = opts.Placeholder
	state.DefaultValue = defaultVal
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

	var value *selectbox.Value
	if state.Value != nil {
		value = &selectbox.Value{
			Value: opts.Options[*state.Value],
			Index: conv.SafeValue(state.Value),
		}
	}

	return value
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
		Label:        state.Label,
		Value:        state.Value,
		Options:      state.Options,
		Placeholder:  state.Placeholder,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
	}
}

func convertSelectboxDataToState(id uuid.UUID, data *websocket.SelectboxData) *selectbox.State {
	if data == nil {
		return nil
	}
	return &selectbox.State{
		ID:           id,
		Label:        data.Label,
		Value:        data.Value,
		Options:      data.Options,
		Placeholder:  data.Placeholder,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
	}
}
