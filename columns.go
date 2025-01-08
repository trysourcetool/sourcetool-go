package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/columns"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

func (b *uiBuilder) Columns(cols int, opts ...columns.Option) []UIBuilder {
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

	columnsOpts := &options.ColumnsOptions{
		Columns: cols,
	}
	for _, option := range opts {
		option.Apply(columnsOpts)
	}

	widgetID := b.generateColumnsID(path)
	weights := columnsOpts.Weight
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

	columnsState := &state.ColumnsState{
		ID:      widgetID,
		Columns: cols,
	}
	sess.State.Set(widgetID, columnsState)

	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
		SessionID:  sess.ID.String(),
		PageID:     page.id.String(),
		WidgetID:   widgetID.String(),
		WidgetType: state.WidgetTypeColumns.String(),
		Path:       path,
		Data:       convertStateToColumnsData(columnsState),
	})

	builders := make([]UIBuilder, cols)
	for i := 0; i < cols; i++ {
		columnCursor := newCursor()
		columnCursor.parentPath = append(path, i)

		columnPath := append(path, i)
		widgetID := b.generateColumnItemID(columnPath)
		columnItemState := &state.ColumnItemState{
			ID:     widgetID,
			Weight: float64(weights[i]) / float64(totalWeight),
		}
		sess.State.Set(widgetID, columnItemState)

		log.Printf("Path: %v\n", columnPath)

		b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodRenderWidget, &websocket.RenderWidgetPayload{
			SessionID:  sess.ID.String(),
			PageID:     page.id.String(),
			WidgetID:   widgetID.String(),
			WidgetType: state.WidgetTypeColumnItem.String(),
			Path:       columnPath,
			Data:       convertStateToColumnItemData(columnItemState),
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
	return uuid.NewV5(page.id, state.WidgetTypeColumns.String()+"-"+path.String())
}

func (b *uiBuilder) generateColumnItemID(path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeColumnItem.String()+"-"+path.String())
}

func convertStateToColumnsData(state *state.ColumnsState) *websocket.ColumnsData {
	return &websocket.ColumnsData{
		Columns: state.Columns,
	}
}

func convertColumnsDataToState(data *websocket.ColumnsData) *state.ColumnsState {
	return &state.ColumnsState{
		Columns: data.Columns,
	}
}

func convertStateToColumnItemData(state *state.ColumnItemState) *websocket.ColumnItemData {
	return &websocket.ColumnItemData{
		Weight: state.Weight,
	}
}

func convertColumnItemDataToState(data *websocket.ColumnItemData) *state.ColumnItemState {
	return &state.ColumnItemState{
		Weight: data.Weight,
	}
}
