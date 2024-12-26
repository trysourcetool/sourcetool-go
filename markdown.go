package sourcetool

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/markdown"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

const widgetTypeMarkdown = "markdown"

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
	path := cursor.getDeltaPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateMarkdownID(body, path)
	state := sess.State.GetMarkdown(widgetID)
	if state == nil {
		// Set initial state
		state = &markdown.State{
			ID:   widgetID,
			Body: opts.Body,
		}
		sess.State.Set(widgetID, state)
	}

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeNumberInput,
		Data:       state,
	})

	cursor.next()
}

func (b *uiBuilder) generateMarkdownID(body string, path []int) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	strPath := make([]string, len(path))
	for i, num := range path {
		strPath[i] = fmt.Sprint(num)
	}
	return uuid.NewV5(page.id, widgetTypeMarkdown+"-"+body+"-"+strings.Join(strPath, ""))
}
