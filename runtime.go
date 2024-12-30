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
	"github.com/trysourcetool/sourcetool-go/internal/checkbox"
	"github.com/trysourcetool/sourcetool-go/internal/columnitem"
	"github.com/trysourcetool/sourcetool-go/internal/columns"
	"github.com/trysourcetool/sourcetool-go/internal/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/form"
	"github.com/trysourcetool/sourcetool-go/internal/markdown"
	"github.com/trysourcetool/sourcetool-go/internal/multiselect"
	"github.com/trysourcetool/sourcetool-go/internal/numberinput"
	"github.com/trysourcetool/sourcetool-go/internal/selectbox"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/table"
	"github.com/trysourcetool/sourcetool-go/internal/textarea"
	"github.com/trysourcetool/sourcetool-go/internal/textinput"
	"github.com/trysourcetool/sourcetool-go/internal/timeinput"
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

	newWidgetStates, err := buildNewWidgetStates(states, sess)
	if err != nil {
		return err
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

func buildNewWidgetStates(states map[uuid.UUID]json.RawMessage, sess *session.Session) (map[uuid.UUID]session.WidgetState, error) {
	widgetStates := make(map[uuid.UUID]session.WidgetState)
	for id, state := range states {
		currentState := sess.State.Get(id)
		if currentState == nil {
			continue
		}
		switch currentState.GetType() {
		case textinput.WidgetType:
			var textInputData websocket.TextInputData
			if err := json.Unmarshal(state, &textInputData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal text input state: %v", err)
			}
			widgetStates[id] = convertTextInputDataToState(id, &textInputData)
		case numberinput.WidgetType:
			var numberInputData websocket.NumberInputData
			if err := json.Unmarshal(state, &numberInputData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal number input state: %v", err)
			}
			widgetStates[id] = convertNumberInputDataToState(&numberInputData)
		case dateinput.WidgetType:
			var dateInputData websocket.DateInputData
			if err := json.Unmarshal(state, &dateInputData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal date input state: %v", err)
			}
			dateInputState, ok := currentState.(*dateinput.State)
			if !ok {
				return nil, fmt.Errorf("invalid date input state: %v", currentState)
			}

			location := dateInputState.Location
			if location == nil {
				location = time.Local
			}

			newState, err := convertDateInputDataToState(id, &dateInputData, location)
			if err != nil {
				return nil, fmt.Errorf("failed to convert date input data: %v", err)
			}

			widgetStates[id] = newState
		case timeinput.WidgetType:
			var timeInputData websocket.TimeInputData
			if err := json.Unmarshal(state, &timeInputData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal time input state: %v", err)
			}
			timeInputState, ok := currentState.(*timeinput.State)
			if !ok {
				return nil, fmt.Errorf("invalid time input state: %v", currentState)
			}

			location := timeInputState.Location
			if location == nil {
				location = time.Local
			}

			newState, err := convertTimeInputDataToState(id, &timeInputData, location)
			if err != nil {
				return nil, fmt.Errorf("failed to convert time input data: %v", err)
			}

			widgetStates[id] = newState
		case form.WidgetType:
			var formData websocket.FormData
			if err := json.Unmarshal(state, &formData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal form state: %v", err)
			}
			widgetStates[id] = convertFormDataToState(&formData)
		case button.WidgetType:
			var buttonData websocket.ButtonData
			if err := json.Unmarshal(state, &buttonData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal button state: %v", err)
			}
			widgetStates[id] = convertButtonDataToState(&buttonData)
		case markdown.WidgetType:
			var markdownData websocket.MarkdownData
			if err := json.Unmarshal(state, &markdownData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal markdown state: %v", err)
			}
			widgetStates[id] = convertMarkdownDataToState(&markdownData)
		case columns.WidgetType:
			var columnsData websocket.ColumnsData
			if err := json.Unmarshal(state, &columnsData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal columns state: %v", err)
			}
			widgetStates[id] = convertColumnsDataToState(&columnsData)
		case columnitem.WidgetType:
			var columnItemData websocket.ColumnItemData
			if err := json.Unmarshal(state, &columnItemData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal column item state: %v", err)
			}
			widgetStates[id] = convertColumnItemDataToState(&columnItemData)
		case table.WidgetType:
			var tableData websocket.TableData
			if err := json.Unmarshal(state, &tableData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal table state: %v", err)
			}
			widgetStates[id] = convertTableDataToState(id, &tableData)
		case selectbox.WidgetType:
			var selectboxData websocket.SelectboxData
			if err := json.Unmarshal(state, &selectboxData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal selectbox state: %v", err)
			}
			widgetStates[id] = convertSelectboxDataToState(id, &selectboxData)
		case multiselect.WidgetType:
			var multiSelectData websocket.MultiSelectData
			if err := json.Unmarshal(state, &multiSelectData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal multiSelect state: %v", err)
			}
			widgetStates[id] = convertMultiSelectDataToState(id, &multiSelectData)
		case checkbox.WidgetType:
			var checkboxData websocket.CheckboxData
			if err := json.Unmarshal(state, &checkboxData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal checkbox state: %v", err)
			}
			widgetStates[id] = convertCheckboxDataToState(id, &checkboxData)
		case textarea.WidgetType:
			var textareaData websocket.TextAreaData
			if err := json.Unmarshal(state, &textareaData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal textarea state: %v", err)
			}
			widgetStates[id] = convertTextAreaDataToState(id, &textareaData)
		default:
			return nil, fmt.Errorf("unknown widget type: %s", currentState.GetType())
		}
	}
	return widgetStates, nil
}
