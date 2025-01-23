package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/selectbox"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
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

	for _, o := range opts {
		o.Apply(selectboxOpts)
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

	selectboxProto := convertStateToSelectboxProto(selectboxState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_Selectbox{
				Selectbox: selectboxProto,
			},
		},
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

func convertStateToSelectboxProto(state *state.SelectboxState) *widgetv1.Selectbox {
	if state == nil {
		return nil
	}
	var value *int64
	if state.Value != nil {
		v := int64(*state.Value)
		value = &v
	}
	var defaultValue *int64
	if state.DefaultValue != nil {
		v := int64(*state.DefaultValue)
		defaultValue = &v
	}
	return &widgetv1.Selectbox{
		Label:        state.Label,
		Value:        value,
		Options:      state.Options,
		Placeholder:  state.Placeholder,
		DefaultValue: defaultValue,
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}

func convertSelectboxProtoToState(id uuid.UUID, data *widgetv1.Selectbox) *state.SelectboxState {
	if data == nil {
		return nil
	}
	var value *int
	if data.Value != nil {
		v := int(*data.Value)
		value = &v
	}
	var defaultValue *int
	if data.DefaultValue != nil {
		v := int(*data.DefaultValue)
		defaultValue = &v
	}
	return &state.SelectboxState{
		ID:           id,
		Label:        data.Label,
		Value:        value,
		Options:      data.Options,
		Placeholder:  data.Placeholder,
		DefaultValue: defaultValue,
		Required:     data.Required,
		Disabled:     data.Disabled,
	}
}
