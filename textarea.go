package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/textarea"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

const widgetTypeTextArea = "textArea"

func (b *uiBuilder) TextArea(label string, options ...textarea.Option) string {
	defaultMinLines := 2
	opts := &textarea.Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: "",
		Required:     false,
		MaxLength:    nil,
		MinLength:    nil,
		MinLines:     &defaultMinLines,
		MaxLines:     nil,
		AutoResize:   true,
	}

	for _, option := range options {
		option(opts)
	}

	sess := b.session
	if sess == nil {
		return ""
	}
	page := b.page
	if page == nil {
		return ""
	}
	cursor := b.cursor
	if cursor == nil {
		return ""
	}
	path := cursor.getPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateTextAreaID(label, path)
	state := sess.State.GetTextArea(widgetID)
	if state == nil {
		// Set initial state
		state = &textarea.State{
			ID:           widgetID,
			Label:        opts.Label,
			Value:        textarea.ReturnValue(opts.DefaultValue),
			Placeholder:  opts.Placeholder,
			DefaultValue: opts.DefaultValue,
			Required:     opts.Required,
			MaxLength:    opts.MaxLength,
			MinLength:    opts.MinLength,
			MaxLines:     opts.MaxLines,
			MinLines:     opts.MinLines,
			AutoResize:   opts.AutoResize,
		}
		sess.State.Set(widgetID, state)
	}
	returnValue := state.Value

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeTextArea,
		Path:       path,
		Data:       state,
	})

	cursor.next()

	return string(returnValue)
}

func (b *uiBuilder) generateTextAreaID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, widgetTypeTextArea+"-"+label+"-"+path.String())
}
