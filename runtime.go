package sourcetool

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/button"
	"github.com/trysourcetool/sourcetool-go/internal/columnitem"
	"github.com/trysourcetool/sourcetool-go/internal/columns"
	"github.com/trysourcetool/sourcetool-go/internal/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/form"
	"github.com/trysourcetool/sourcetool-go/internal/markdown"
	"github.com/trysourcetool/sourcetool-go/internal/numberinput"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/table"
	"github.com/trysourcetool/sourcetool-go/internal/textarea"
	"github.com/trysourcetool/sourcetool-go/internal/textinput"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

var once sync.Once

type runtime struct {
	wsClient       websocket.Client
	sessionManager *session.SessionManager
	pageManager    *pageManager
}

func startRuntime(apiKey, endpoint string, pages map[uuid.UUID]*page) *runtime {
	var r *runtime
	once.Do(func() {
		r = &runtime{
			sessionManager: session.NewSessionManager(),
			pageManager:    newPageManager(pages),
		}

		wsClient, err := websocket.NewClient(websocket.Config{
			URL:            endpoint,
			APIKey:         apiKey,
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
			log.Fatalf("failed to create websocket client: %v", err)
		}

		r.wsClient = wsClient
		wsClient.RegisterHandler(websocket.MessageMethodInitializeClient, r.handleInitializeCilent)
		wsClient.RegisterHandler(websocket.MessageMethodRerunPage, r.handleRerunPage)
		wsClient.RegisterHandler(websocket.MessageMethodCloseSession, r.handleCloseSession)

		r.sendInitializeHost(apiKey, pages)
	})

	return r
}

func (r *runtime) sendInitializeHost(apiKey string, pages map[uuid.UUID]*page) {
	pagesPayload := make([]*websocket.InitializeHostPagePayload, 0, len(pages))
	for _, page := range pages {
		pagesPayload = append(pagesPayload, &websocket.InitializeHostPagePayload{
			ID:   page.id.String(),
			Name: page.name,
		})
	}

	resp, err := r.wsClient.EnqueueWithResponse(uuid.Must(uuid.NewV4()).String(), websocket.MessageMethodInitializeHost, websocket.InitializeHostPayload{
		APIKey:     apiKey,
		SDKName:    "sourcetool-go",
		SDKVersion: "0.1.0",
		Pages:      pagesPayload,
	})
	if err != nil {
		log.Fatalf("failed to send initialize host message: %v", err)
	}
	if resp.Error != nil {
		log.Fatalf("initialize host message failed: %v", resp.Error)
	}

	log.Printf("initialize host message sent: %v", resp)
}

func (r *runtime) handleInitializeCilent(msg *websocket.Message) error {
	var p websocket.InitializeClientPayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	sessionID, err := uuid.FromString(p.SessionID)
	if err != nil {
		return err
	}
	pageID, err := uuid.FromString(p.PageID)
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

func (r *runtime) handleRerunPage(msg *websocket.Message) error {
	var p websocket.RerunPagePayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	sessionID, err := uuid.FromString(p.SessionID)
	if err != nil {
		return err
	}
	sess := r.sessionManager.GetSession(sessionID)
	if sess == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	pageID, err := uuid.FromString(p.PageID)
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

	var states map[uuid.UUID]json.RawMessage
	if err := json.Unmarshal(p.State, &states); err != nil {
		return fmt.Errorf("failed to unmarshal state: %v", err)
	}

	widgetStates := make(map[uuid.UUID]session.WidgetState)
	for id, state := range states {
		currentState := sess.State.Get(id)
		if currentState == nil {
			continue
		}
		switch currentState.GetType() {
		case textinput.WidgetType:
			var textInputState textinput.State
			if err := json.Unmarshal(state, &textInputState); err != nil {
				return fmt.Errorf("failed to unmarshal text input state: %v", err)
			}
			widgetStates[id] = &textInputState
		case numberinput.WidgetType:
			var numberInputState numberinput.State
			if err := json.Unmarshal(state, &numberInputState); err != nil {
				return fmt.Errorf("failed to unmarshal number input state: %v", err)
			}
			widgetStates[id] = &numberInputState
		case dateinput.WidgetType:
			var dateInputData websocket.DateInputData
			if err := json.Unmarshal(state, &dateInputData); err != nil {
				return fmt.Errorf("failed to unmarshal date input state: %v", err)
			}
			dateInputState, ok := currentState.(*dateinput.State)
			if !ok {
				return fmt.Errorf("invalid date input state: %v", currentState)
			}
			location := time.Local
			if dateInputState.Location != nil {
				location = dateInputState.Location
			}
			var val, defaultVal, maxVal, minVal time.Time
			if dateInputData.Value != "" {
				val, err = time.ParseInLocation(time.DateOnly, dateInputData.Value, location)
				if err != nil {
					return err
				}
			}
			if dateInputData.DefaultValue != "" {
				defaultVal, err = time.ParseInLocation(time.DateOnly, dateInputData.DefaultValue, location)
				if err != nil {
					return err
				}
			}
			if dateInputData.MaxValue != "" {
				maxVal, err = time.ParseInLocation(time.DateOnly, dateInputData.MaxValue, location)
				if err != nil {
					return err
				}
			}
			if dateInputData.MinValue != "" {
				minVal, err = time.ParseInLocation(time.DateOnly, dateInputData.MinValue, location)
				if err != nil {
					return err
				}
			}
			widgetStates[id] = &dateinput.State{
				ID:           id,
				Value:        &val,
				Label:        dateInputState.Label,
				Placeholder:  dateInputState.Placeholder,
				DefaultValue: &defaultVal,
				Required:     dateInputState.Required,
				Format:       dateInputState.Format,
				MaxValue:     &maxVal,
				MinValue:     &minVal,
			}
		case form.WidgetType:
			var formState form.State
			if err := json.Unmarshal(state, &formState); err != nil {
				return fmt.Errorf("failed to unmarshal form state: %v", err)
			}
			widgetStates[id] = &formState
		case button.WidgetType:
			var buttonState button.State
			if err := json.Unmarshal(state, &buttonState); err != nil {
				return fmt.Errorf("failed to unmarshal button state: %v", err)
			}
			widgetStates[id] = &buttonState
		case markdown.WidgetType:
			var markdownState markdown.State
			if err := json.Unmarshal(state, &markdownState); err != nil {
				return fmt.Errorf("failed to unmarshal markdown state: %v", err)
			}
			widgetStates[id] = &markdownState
		case columns.WidgetType:
			var columnsState columns.State
			if err := json.Unmarshal(state, &columnsState); err != nil {
				return fmt.Errorf("failed to unmarshal columns state: %v", err)
			}
			widgetStates[id] = &columnsState
		case columnitem.WidgetType:
			var columnItemState columnitem.State
			if err := json.Unmarshal(state, &columnItemState); err != nil {
				return fmt.Errorf("failed to unmarshal column item state: %v", err)
			}
			widgetStates[id] = &columnItemState
		case table.WidgetType:
			var tableState table.State
			if err := json.Unmarshal(state, &tableState); err != nil {
				return fmt.Errorf("failed to unmarshal table state: %v", err)
			}
			widgetStates[id] = &tableState
		case textarea.WidgetType:
			var textareaState textarea.State
			if err := json.Unmarshal(state, &textareaState); err != nil {
				return fmt.Errorf("failed to unmarshal textarea state: %v", err)
			}
			widgetStates[id] = &textareaState
		default:
			return fmt.Errorf("unknown widget type: %s", currentState.GetType())
		}
	}

	sess.State.SetStates(widgetStates)

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

func (r *runtime) handleCloseSession(msg *websocket.Message) error {
	var p websocket.CloseSessionPayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	sessionID, err := uuid.FromString(p.SessionID)
	if err != nil {
		return err
	}

	r.sessionManager.DeleteSession(sessionID)

	return nil
}
