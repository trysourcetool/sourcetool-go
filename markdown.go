package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) Markdown(body string) {
	markdownOpts := &options.MarkdownOptions{
		Body: body,
	}

	sess := b.session
	if sess == nil {
		return
	}
	page := b.page
	if page == nil {
		return
	}
	cursor := b.cursor
	if cursor == nil {
		return
	}
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateMarkdownID(body, path)
	markdownState := sess.State.GetMarkdown(widgetID)
	if markdownState == nil {
		markdownState = &state.MarkdownState{
			ID: widgetID,
		}
	}
	markdownState.Body = markdownOpts.Body
	sess.State.Set(widgetID, markdownState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeMarkdown.String(),
		Path:       path,
		Data:       convertStateToMarkdownData(markdownState),
	})

	cursor.next()
}

func (b *uiBuilder) generateMarkdownID(body string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeMarkdown.String()+"-"+body+"-"+path.String())
}

func convertStateToMarkdownData(state *state.MarkdownState) *websocket.MarkdownData {
	return &websocket.MarkdownData{
		Body: state.Body,
	}
}

func convertMarkdownDataToState(data *websocket.MarkdownData) *state.MarkdownState {
	return &state.MarkdownState{
		Body: data.Body,
	}
}
