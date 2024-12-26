package sourcetool

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/numberinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

const widgetTypeNumberInput = "numberInput"

func (b *uiBuilder) NumberInput(label string, options ...numberinput.Option) float64 {
	opts := &numberinput.Options{
		Label:        label,
		Placeholder:  "",
		DefaultValue: 0,
		Required:     false,
		MaxValue:     nil,
		MinValue:     nil,
	}

	for _, option := range options {
		option(opts)
	}

	sess := b.session
	if sess == nil {
		return 0
	}
	page := b.page
	if page == nil {
		return 0
	}
	cursor := b.cursor
	if cursor == nil {
		return 0
	}
	path := cursor.getDeltaPath()

	log.Printf("Session ID: %s", sess.ID.String())
	log.Printf("Page ID: %s", page.id.String())
	log.Printf("Path: %v\n", path)

	widgetID := b.generateNumberInputID(label, path)
	state := sess.State.GetNumberInput(widgetID)
	if state == nil {
		// Set initial state
		state = &numberinput.State{
			ID:           widgetID,
			Label:        opts.Label,
			Value:        numberinput.ReturnValue(opts.DefaultValue),
			Placeholder:  opts.Placeholder,
			DefaultValue: opts.DefaultValue,
			Required:     opts.Required,
			MaxValue:     opts.MaxValue,
			MinValue:     opts.MinValue,
		}
		sess.State.Set(widgetID, state)
	}
	returnValue := state.Value

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: widgetTypeNumberInput,
		Data:       state,
	})

	cursor.next()

	return float64(returnValue)
}

func (b *uiBuilder) generateNumberInputID(label string, path []int) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	strPath := make([]string, len(path))
	for i, num := range path {
		strPath[i] = fmt.Sprint(num)
	}
	return uuid.NewV5(page.id, widgetTypeNumberInput+"-"+label+"-"+strings.Join(strPath, ""))
}
