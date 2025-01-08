package sourcetool

import (
	"github.com/gofrs/uuid/v5"
	"github.com/trysourcetool/sourcetool-go/form"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) Form(buttonLabel string, opts ...form.Option) (UIBuilder, bool) {
	formOpts := &options.FormOptions{
		ButtonLabel:    buttonLabel,
		ButtonDisabled: false,
		ClearOnSubmit:  false,
	}

	for _, o := range opts {
		o.Apply(formOpts)
	}

	sess := b.session
	if sess == nil {
		return b, false
	}
	page := b.page
	if page == nil {
		return b, false
	}
	cursor := b.cursor
	if cursor == nil {
		return b, false
	}
	path := cursor.getPath()

	widgetID := b.generateFormID(path)
	formState := sess.State.GetForm(widgetID)
	if formState == nil {
		formState = &state.FormState{
			ID:    widgetID,
			Value: false,
		}
	}
	formState.ButtonLabel = formOpts.ButtonLabel
	formState.ButtonDisabled = formOpts.ButtonDisabled
	formState.ClearOnSubmit = formOpts.ClearOnSubmit
	sess.State.Set(widgetID, formState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeForm.String(),
		Path:       path,
		Data:       convertStateToFormData(formState),
	})

	cursor.next()

	childCursor := newCursor()
	childCursor.parentPath = path

	childBuilder := &uiBuilder{
		runtime: b.runtime,
		session: sess,
		page:    page,
		cursor:  childCursor,
	}

	return childBuilder, formState.Value
}

func (b *uiBuilder) generateFormID(path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeForm.String()+"-"+path.String())
}

func convertStateToFormData(state *state.FormState) *websocket.FormData {
	if state == nil {
		return nil
	}
	return &websocket.FormData{
		Value:          state.Value,
		ButtonLabel:    state.ButtonLabel,
		ButtonDisabled: state.ButtonDisabled,
		ClearOnSubmit:  state.ClearOnSubmit,
	}
}

func convertFormDataToState(data *websocket.FormData) *state.FormState {
	if data == nil {
		return nil
	}
	return &state.FormState{
		Value:          data.Value,
		ButtonLabel:    data.ButtonLabel,
		ButtonDisabled: data.ButtonDisabled,
		ClearOnSubmit:  data.ClearOnSubmit,
	}
}
