package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/columnitem"
	"github.com/trysourcetool/sourcetool-go/internal/columns"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) Columns(cols int, options ...columns.Option) []UIBuilder {
	if cols < 1 {
		return nil
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

	log.Printf("Path: %v\n", path)

	opts := &columns.Options{
		Columns: cols,
	}
	for _, option := range options {
		option(opts)
	}

	widgetID := b.generateColumnsID(path)
	weights := opts.Weight
	if len(weights) == 0 || len(weights) != cols {
		weights = make([]int, cols)
		for i := range weights {
			weights[i] = 1
		}
	}

	for _, w := range weights {
		if w <= 0 {
			weights = make([]int, cols)
			for i := range weights {
				weights[i] = 1
			}
			break
		}
	}

	totalWeight := 0
	for _, w := range weights {
		totalWeight += w
	}

	columnsState := &columns.State{
		ID:      widgetID,
		Columns: cols,
	}
	sess.State.Set(widgetID, columnsState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: columns.WidgetType,
		Path:       path,
		Data:       columnsState,
	})

	builders := make([]UIBuilder, cols)
	for i := 0; i < cols; i++ {
		columnCursor := newCursor()
		columnCursor.parentPath = append(path, i)

		columnPath := append(path, i)
		widgetID := b.generateColumnItemID(columnPath)
		columnItemState := &columnitem.State{
			ID:     widgetID,
			Weight: float64(weights[i]) / float64(totalWeight),
		}
		sess.State.Set(widgetID, columnItemState)

		log.Printf("Path: %v\n", columnPath)

		b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
			SessionID:  sess.ID.String(),
			PageID:     page.id.String(),
			WidgetID:   widgetID.String(),
			WidgetType: columnitem.WidgetType,
			Path:       columnPath,
			Data:       columnItemState,
		})

		builders[i] = &uiBuilder{
			runtime: b.runtime,
			context: b.context,
			cursor:  columnCursor,
			session: sess,
			page:    page,
		}
	}

	cursor.next()

	return builders
}

func (b *uiBuilder) generateColumnsID(path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, columns.WidgetType+"-"+path.String())
}

func (b *uiBuilder) generateColumnItemID(path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, columnitem.WidgetType+"-"+path.String())
}
