package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/checkboxgroup"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) CheckboxGroup(label string, options ...checkboxgroup.Option) *checkboxgroup.Value {
	opts := &checkboxgroup.Options{
		Label:        label,
		DefaultValue: nil,
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

	widgetID := b.generateCheckboxGroupID(label, path)
	state := sess.State.GetCheckboxGroup(widgetID)
	if state == nil {
		state = &checkboxgroup.State{
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
	state.DefaultValue = defaultVal
	state.Required = opts.Required
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: checkboxgroup.WidgetType,
		Path:       path,
		Data:       convertStateToCheckboxGroupData(state),
	})

	cursor.next()

	var value *checkboxgroup.Value
	if state.Value != nil {
		value = &checkboxgroup.Value{
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

func (b *uiBuilder) generateCheckboxGroupID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, checkboxgroup.WidgetType+"-"+label+"-"+path.String())
}

func convertStateToCheckboxGroupData(state *checkboxgroup.State) *websocket.CheckboxGroupData {
	if state == nil {
		return nil
	}
	return &websocket.CheckboxGroupData{
		Label:        state.Label,
		Value:        state.Value,
		Options:      state.Options,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
	}
}

func convertCheckboxGroupDataToState(id uuid.UUID, data *websocket.CheckboxGroupData) *checkboxgroup.State {
	if data == nil {
		return nil
	}
	return &checkboxgroup.State{
		ID:           id,
		Label:        data.Label,
		Value:        data.Value,
		Options:      data.Options,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
	}
}
