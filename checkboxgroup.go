package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/checkboxgroup"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) CheckboxGroup(label string, opts ...checkboxgroup.Option) *checkboxgroup.Value {
	checkboxGroupOpts := &options.CheckboxGroupOptions{
		Label:        label,
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		FormatFunc:   nil,
	}

	for _, o := range opts {
		o.Apply(checkboxGroupOpts)
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
	if len(checkboxGroupOpts.DefaultValue) != 0 {
		for _, o := range checkboxGroupOpts.DefaultValue {
			for i, opt := range checkboxGroupOpts.Options {
				if opt == o {
					defaultVal = append(defaultVal, i)
					break
				}
			}
		}
	}

	widgetID := b.generateCheckboxGroupID(label, path)
	checkboxGroupState := sess.State.GetCheckboxGroup(widgetID)
	if checkboxGroupState == nil {
		checkboxGroupState = &state.CheckboxGroupState{
			ID:    widgetID,
			Value: defaultVal,
		}
	}
	if checkboxGroupOpts.FormatFunc == nil {
		checkboxGroupOpts.FormatFunc = func(v string, i int) string {
			return v
		}
	}

	displayVals := make([]string, len(checkboxGroupOpts.Options))
	for i, v := range checkboxGroupOpts.Options {
		displayVals[i] = checkboxGroupOpts.FormatFunc(v, i)
	}

	checkboxGroupState.Label = checkboxGroupOpts.Label
	checkboxGroupState.Options = displayVals
	checkboxGroupState.DefaultValue = defaultVal
	checkboxGroupState.Required = checkboxGroupOpts.Required
	checkboxGroupState.Disabled = checkboxGroupOpts.Disabled
	sess.State.Set(widgetID, checkboxGroupState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeCheckbox.String(),
		Path:       path,
		Data:       convertStateToCheckboxGroupData(checkboxGroupState),
	})

	cursor.next()

	var value *checkboxgroup.Value
	if checkboxGroupState.Value != nil {
		value = &checkboxgroup.Value{
			Values:  make([]string, len(checkboxGroupState.Value)),
			Indexes: make([]int, len(checkboxGroupState.Value)),
		}
		for i, idx := range checkboxGroupState.Value {
			value.Values[i] = checkboxGroupOpts.Options[idx]
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
	return uuid.NewV5(page.id, state.WidgetTypeCheckboxGroup.String()+"-"+label+"-"+path.String())
}

func convertStateToCheckboxGroupData(state *state.CheckboxGroupState) *websocket.CheckboxGroupData {
	if state == nil {
		return nil
	}
	return &websocket.CheckboxGroupData{
		Label:        state.Label,
		Value:        state.Value,
		Options:      state.Options,
		DefaultValue: state.DefaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}

func convertCheckboxGroupDataToState(id uuid.UUID, data *websocket.CheckboxGroupData) *state.CheckboxGroupState {
	if data == nil {
		return nil
	}
	return &state.CheckboxGroupState{
		ID:           id,
		Label:        data.Label,
		Value:        data.Value,
		Options:      data.Options,
		DefaultValue: data.DefaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
