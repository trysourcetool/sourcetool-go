package sourcetool

import (
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/columns"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
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
	for _, o := range opts {
		o.Apply(columnsOpts)
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

	columns := convertStateToColumnsProto(columnsState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_Columns{
				Columns: columns,
			},
		},
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

		columnItem := convertStateToColumnItemProto(columnItemState)
		b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
			SessionId: sess.ID.String(),
			PageId:    page.id.String(),
			Path:      convertPathToInt32Slice(path),
			Widget: &widgetv1.Widget{
				Id: widgetID.String(),
				Type: &widgetv1.Widget_ColumnItem{
					ColumnItem: columnItem,
				},
			},
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

func convertStateToColumnsProto(state *state.ColumnsState) *widgetv1.Columns {
	return &widgetv1.Columns{
		Columns: int32(state.Columns),
	}
}

func convertColumnsProtoToState(id uuid.UUID, data *widgetv1.Columns) *state.ColumnsState {
	return &state.ColumnsState{
		ID:      id,
		Columns: int(data.Columns),
	}
}

func convertStateToColumnItemProto(state *state.ColumnItemState) *widgetv1.ColumnItem {
	return &widgetv1.ColumnItem{
		Weight: state.Weight,
	}
}

func convertColumnItemProtoToState(id uuid.UUID, data *widgetv1.ColumnItem) *state.ColumnItemState {
	return &state.ColumnItemState{
		ID:     id,
		Weight: data.Weight,
	}
}
