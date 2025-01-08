package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	"github.com/trysourcetool/sourcetool-go/selectbox"
)

func (b *uiBuilder) Selectbox(label string, opts ...selectbox.Option) *selectbox.Value {
	selectboxOpts := &options.SelectboxOptions{
		Label:        label,
		DefaultValue: nil,
		Placeholder:  "",
		Required:     false,
		Disabled:     false,
		FormatFunc:   nil,
	}

	for _, option := range opts {
		option.Apply(selectboxOpts)
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
	if selectboxOpts.DefaultValue != nil {
		for i, o := range selectboxOpts.Options {
			if conv.SafeValue(selectboxOpts.DefaultValue) == o {
				defaultVal = &i
				break
			}
		}
	}

	widgetID := b.generateSelectboxID(label, path)
	selectboxState := sess.State.GetSelectbox(widgetID)
	if selectboxState == nil {
		selectboxState = &state.SelectboxState{
			ID:           widgetID,
			Value:        defaultVal,
			DefaultValue: defaultVal,
		}
	}

	if selectboxOpts.FormatFunc == nil {
		selectboxOpts.FormatFunc = func(v string, i int) string {
			return v
		}
	}

	displayVals := make([]string, len(selectboxOpts.Options))
	for i, v := range selectboxOpts.Options {
		displayVals[i] = selectboxOpts.FormatFunc(v, i)
	}

	selectboxState.Label = selectboxOpts.Label
	selectboxState.Options = displayVals
	selectboxState.Placeholder = selectboxOpts.Placeholder
	selectboxState.DefaultValue = defaultVal
	selectboxState.Required = selectboxOpts.Required
	selectboxState.Disabled = selectboxOpts.Disabled
	sess.State.Set(widgetID, selectboxState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeSelectbox.String(),
		Path:       path,
		Data:       convertStateToSelectboxData(selectboxState),
	})

	cursor.next()

	var value *selectbox.Value
	if selectboxState.Value != nil {
		value = &selectbox.Value{
			Value: selectboxOpts.Options[*selectboxState.Value],
			Index: conv.SafeValue(selectboxState.Value),
		}
	}

	return value
}

func (b *uiBuilder) generateSelectboxID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeSelectbox.String()+"-"+label+"-"+path.String())
}

func convertStateToSelectboxData(state *state.SelectboxState) *websocket.SelectboxData {
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
		Disabled:     state.Disabled,
	}
}

func convertSelectboxDataToState(id uuid.UUID, data *websocket.SelectboxData) *state.SelectboxState {
	if data == nil {
		return nil
	}
	return &state.SelectboxState{
		ID:           id,
		Label:        data.Label,
		Value:        data.Value,
		Options:      data.Options,
		Placeholder:  data.Placeholder,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
