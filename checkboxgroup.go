package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/checkboxgroup"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
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

	checkboxGroupProto := convertStateToCheckboxGroupProto(checkboxGroupState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_CheckboxGroup{
				CheckboxGroup: checkboxGroupProto,
			},
		},
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

func convertStateToCheckboxGroupProto(state *state.CheckboxGroupState) *widgetv1.CheckboxGroup {
	if state == nil {
		return nil
	}
	value := make([]int64, len(state.Value))
	if len(state.Value) != 0 {
		for i, v := range state.Value {
			value[i] = int64(v)
		}
	}
	defaultValue := make([]int64, len(state.DefaultValue))
	if len(state.DefaultValue) != 0 {
		for i, v := range state.DefaultValue {
			defaultValue[i] = int64(v)
		}
	}
	return &widgetv1.CheckboxGroup{
		Label:        state.Label,
		Value:        value,
		Options:      state.Options,
		DefaultValue: defaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}

func convertCheckboxGroupProtoToState(id uuid.UUID, data *widgetv1.CheckboxGroup) *state.CheckboxGroupState {
	if data == nil {
		return nil
	}
	value := make([]int, len(data.Value))
	if len(data.Value) != 0 {
		for i, v := range data.Value {
			value[i] = int(v)
		}
	}
	defaultValue := make([]int, len(data.DefaultValue))
	if len(data.DefaultValue) != 0 {
		for i, v := range data.DefaultValue {
			defaultValue[i] = int(v)
		}
	}
	return &state.CheckboxGroupState{
		ID:           id,
		Label:        data.Label,
		Value:        value,
		Options:      data.Options,
		DefaultValue: defaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
