package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/markdown"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) Markdown(body string, options ...markdown.Option) {
	opts := &markdown.Options{
		Body: body,
	}

	for _, option := range options {
		option(opts)
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
	state := sess.State.GetMarkdown(widgetID)
	if state == nil {
		state = &markdown.State{
			ID: widgetID,
		}
	}
	state.Body = opts.Body
	sess.State.Set(widgetID, state)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: markdown.WidgetType,
		Path:       path,
		Data:       convertStateToMarkdownData(state),
	})

	cursor.next()
}

func (b *uiBuilder) generateMarkdownID(body string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, markdown.WidgetType+"-"+body+"-"+path.String())
}

func convertStateToMarkdownData(state *markdown.State) *websocket.MarkdownData {
	return &websocket.MarkdownData{
		Body: state.Body,
	}
}

func convertMarkdownDataToState(data *websocket.MarkdownData) *markdown.State {
	return &markdown.State{
		Body: data.Body,
	}
}
