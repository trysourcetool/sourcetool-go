package sourcetool

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
	pagev1 "github.com/trysourcetool/sourcetool-proto/go/page/v1"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
)

type runtime struct {
	wsClient       websocket.Client
	sessionManager *session.SessionManager
	pageManager    *pageManager
}

func startRuntime(apiKey, endpoint string, pages map[uuid.UUID]*page) (*runtime, error) {
	r := &runtime{
		sessionManager: session.NewSessionManager(),
		pageManager:    newPageManager(pages),
	}

	wsClient, err := websocket.NewClient(websocket.Config{
		URL:            endpoint,
		APIKey:         apiKey,
		InstanceID:     uuid.Must(uuid.NewV4()),
		PingInterval:   1 * time.Second,
		ReconnectDelay: 1 * time.Second,
		OnReconnecting: func() {
			log.Println("Reconnecting...")
		},
		OnReconnected: func() {
			log.Println("Reconnected!")
			r.sendInitializeHost(apiKey, pages)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket client: %v", err)
	}

	r.wsClient = wsClient
	wsClient.RegisterHandler(func(msg *websocketv1.Message) error {
		switch t := msg.Type.(type) {
		case *websocketv1.Message_InitializeClient:
			return r.handleInitializeClient(t.InitializeClient)
		case *websocketv1.Message_RerunPage:
			return r.handleRerunPage(t.RerunPage)
		case *websocketv1.Message_CloseSession:
			return r.handleCloseSession(t.CloseSession)
		default:
			return fmt.Errorf("unknown message type: %T", t)
		}
	})

	r.sendInitializeHost(apiKey, pages)

	return r, nil
}

func (r *runtime) sendInitializeHost(apiKey string, pages map[uuid.UUID]*page) {
	pagesPayload := make([]*pagev1.Page, 0, len(pages))
	for _, page := range pages {
		pagesPayload = append(pagesPayload, &pagev1.Page{
			Id:     page.id.String(),
			Name:   page.name,
			Route:  page.route,
			Path:   conv.PathToInt32Slice(page.path),
			Groups: page.accessGroups,
		})
	}

	msg := &websocketv1.InitializeHost{
		ApiKey:     apiKey,
		SdkName:    "sourcetool-go",
		SdkVersion: "0.1.0",
		Pages:      pagesPayload,
	}

	resp, err := r.wsClient.EnqueueWithResponse(uuid.Must(uuid.NewV4()).String(), msg)
	if err != nil {
		log.Fatalf("failed to send initialize host message: %v", err)
	}

	if e := resp.GetException(); e != nil {
		log.Fatalf("initialize host message failed: %v", e.Message)
	}

	log.Printf("initialize host message sent: %v", resp)
}

func (r *runtime) handleInitializeClient(msg *websocketv1.InitializeClient) error {
	sessionID, err := uuid.FromString(msg.SessionId)
	if err != nil {
		return err
	}
	pageID, err := uuid.FromString(msg.PageId)
	if err != nil {
		return err
	}

	log.Println("Creating new session with ID:", sessionID)
	session := session.New(sessionID, pageID)
	r.sessionManager.SetSession(session)

	page := r.pageManager.getPage(pageID)
	if page == nil {
		return fmt.Errorf("page not found: %s", pageID)
	}

	ui := &uiBuilder{
		context: context.Background(),
		runtime: r,
		session: session,
		page:    page,
		cursor:  newCursor(),
	}

	if err := page.run(ui); err != nil {
		return fmt.Errorf("failed to run page: %v", err)
	}

	return nil
}

func (r *runtime) handleRerunPage(msg *websocketv1.RerunPage) error {
	sessionID, err := uuid.FromString(msg.SessionId)
	if err != nil {
		return err
	}
	sess := r.sessionManager.GetSession(sessionID)
	if sess == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	pageID, err := uuid.FromString(msg.PageId)
	if err != nil {
		return err
	}
	page := r.pageManager.getPage(pageID)
	if page == nil {
		return fmt.Errorf("page not found: %s", pageID)
	}

	if sess.PageID != pageID {
		sess.State.ResetStates()
	}

	newWidgetStates := make(map[uuid.UUID]session.WidgetState)
	for _, widget := range msg.States {
		switch t := widget.Type.(type) {
		case *widgetv1.Widget_TextInput:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertTextInputProtoToState(id, t.TextInput)
		case *widgetv1.Widget_NumberInput:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertNumberInputProtoToState(id, t.NumberInput)
		case *widgetv1.Widget_DateInput:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			state, err := convertDateInputProtoToState(id, t.DateInput, time.Local)
			if err != nil {
				return err
			}
			newWidgetStates[id] = state
		case *widgetv1.Widget_DateTimeInput:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			state, err := convertDateTimeInputProtoToState(id, t.DateTimeInput, time.Local)
			if err != nil {
				return err
			}
			newWidgetStates[id] = state
		case *widgetv1.Widget_TimeInput:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			state, err := convertTimeInputProtoToState(id, t.TimeInput, time.Local)
			if err != nil {
				return err
			}
			newWidgetStates[id] = state
		case *widgetv1.Widget_Form:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertFormProtoToState(id, t.Form)
		case *widgetv1.Widget_Button:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertButtonProtoToState(id, t.Button)
		case *widgetv1.Widget_Markdown:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertMarkdownProtoToState(id, t.Markdown)
		case *widgetv1.Widget_Columns:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertColumnsProtoToState(id, t.Columns)
		case *widgetv1.Widget_ColumnItem:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertColumnItemProtoToState(id, t.ColumnItem)
		case *widgetv1.Widget_Table:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertTableProtoToState(id, t.Table)
		case *widgetv1.Widget_Selectbox:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertSelectboxProtoToState(id, t.Selectbox)
		case *widgetv1.Widget_MultiSelect:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertMultiSelectProtoToState(id, t.MultiSelect)
		case *widgetv1.Widget_Checkbox:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertCheckboxProtoToState(id, t.Checkbox)
		case *widgetv1.Widget_CheckboxGroup:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertCheckboxGroupProtoToState(id, t.CheckboxGroup)
		case *widgetv1.Widget_Radio:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertRadioProtoToState(id, t.Radio)
		case *widgetv1.Widget_TextArea:
			id, err := uuid.FromString(widget.Id)
			if err != nil {
				return err
			}
			newWidgetStates[id] = convertTextAreaProtoToState(id, t.TextArea)
		default:
			return fmt.Errorf("unknown widget type: %T", t)
		}
	}

	sess.State.SetStates(newWidgetStates)

	ui := &uiBuilder{
		context: context.Background(),
		runtime: r,
		session: sess,
		page:    page,
		cursor:  newCursor(),
	}

	if err := page.run(ui); err != nil {
		return fmt.Errorf("failed to run page: %v", err)
	}

	sess.State.ResetButtons()

	return nil
}

func (r *runtime) handleCloseSession(msg *websocketv1.CloseSession) error {
	sessionID, err := uuid.FromString(msg.SessionId)
	if err != nil {
		return err
	}

	r.sessionManager.DeleteSession(sessionID)

	return nil
}
