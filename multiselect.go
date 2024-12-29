package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/multiselect"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) MultiSelect(label string, options ...multiselect.Option) *multiselect.Value {
	opts := &multiselect.Options{
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

	widgetID := b.generateMultiSelectID(label, path)
	state := sess.State.GetMultiSelect(widgetID)
	if state == nil {
		var defaultVal *multiselect.Value
		if len(opts.DefaultValue) != 0 {
			defaultIndexes := make([]int, len(opts.DefaultValue))
			for i, o := range opts.Options {
				for j, v := range opts.DefaultValue {
					if o == v {
						defaultIndexes[j] = i
					}
				}
			}
			defaultVal = &multiselect.Value{
				Values:  opts.DefaultValue,
				Indexes: defaultIndexes,
			}
		}
		state = &multiselect.State{
			ID:    widgetID,
			Value: defaultVal,
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
	state.DefaultValue = opts.DefaultValue
	state.Required = opts.Required
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: multiselect.WidgetType,
		Path:       path,
		Data:       convertStateToMultiSelectData(state),
	})

	cursor.next()

	return state.Value
}

func (b *uiBuilder) generateMultiSelectID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, multiselect.WidgetType+"-"+label+"-"+path.String())
}

func convertStateToMultiSelectData(state *multiselect.State) *websocket.MultiSelectData {
	if state == nil {
		return nil
	}
	var value *websocket.MultiSelectDataValue
	if state.Value != nil {
		value = &websocket.MultiSelectDataValue{
			Values:  state.Value.Values,
			Indexes: state.Value.Indexes,
		}
	}
	return &websocket.MultiSelectData{
		Label:        state.Label,
		Value:        value,
		Options:      state.Options,
		Placeholder:  state.Placeholder,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
	}
}

func convertMultiSelectDataToState(id uuid.UUID, data *websocket.MultiSelectData) *multiselect.State {
	if data == nil {
		return nil
	}
	var value *multiselect.Value
	if data.Value != nil {
		value = &multiselect.Value{
			Values:  data.Value.Values,
			Indexes: data.Value.Indexes,
		}
	}
	return &multiselect.State{
		ID:           id,
		Label:        data.Label,
		Value:        value,
		Options:      data.Options,
		Placeholder:  data.Placeholder,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
	}
}
