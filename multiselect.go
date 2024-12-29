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

	var defaultVal []int
	if len(opts.DefaultValue) != 0 {
		for _, o := range opts.DefaultValue {
			for i, opt := range opts.Options {
				if opt == o {
					defaultVal = append(defaultVal, i)
					break
				}
			}
		}
	}

	widgetID := b.generateMultiSelectID(label, path)
	state := sess.State.GetMultiSelect(widgetID)
	if state == nil {
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
	state.DefaultValue = defaultVal
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

	var value *multiselect.Value
	if state.Value != nil {
		value = &multiselect.Value{
			Values:  make([]string, len(state.Value)),
			Indexes: make([]int, len(state.Value)),
		}
		for i, idx := range state.Value {
			value.Values[i] = opts.Options[idx]
			value.Indexes[i] = idx
		}
	}

	return value
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
	return &websocket.MultiSelectData{
		Label:        state.Label,
		Value:        state.Value,
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
	return &multiselect.State{
		ID:           id,
		Label:        data.Label,
		Value:        data.Value,
		Options:      data.Options,
		Placeholder:  data.Placeholder,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
	}
}
