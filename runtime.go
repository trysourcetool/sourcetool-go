package sourcetool

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket"
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
	wsClient.RegisterHandler(websocket.MessageMethodInitializeClient, r.handleInitializeCilent)
	wsClient.RegisterHandler(websocket.MessageMethodRerunPage, r.handleRerunPage)
	wsClient.RegisterHandler(websocket.MessageMethodCloseSession, r.handleCloseSession)

	r.sendInitializeHost(apiKey, pages)

	return r, nil
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
	for id, s := range states {
		currentState := sess.State.Get(id)
		if currentState == nil {
			continue
		}
		switch currentState.GetType() {
		case state.WidgetTypeTextInput:
			var textInputData websocket.TextInputData
			if err := json.Unmarshal(s, &textInputData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal text input state: %v", err)
			}
			widgetStates[id] = convertTextInputDataToState(id, &textInputData)
		case state.WidgetTypeNumberInput:
			var numberInputData websocket.NumberInputData
			if err := json.Unmarshal(s, &numberInputData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal number input state: %v", err)
			}
			widgetStates[id] = convertNumberInputDataToState(&numberInputData)
		case state.WidgetTypeDateInput:
			var dateInputData websocket.DateInputData
			if err := json.Unmarshal(s, &dateInputData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal date input state: %v", err)
			}
			dateInputState, ok := currentState.(*state.DateInputState)
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
		case state.WidgetTypeDateTimeInput:
			var dateTimeInputData websocket.DateTimeInputData
			if err := json.Unmarshal(s, &dateTimeInputData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal datetime input state: %v", err)
			}
			dateTimeInputState, ok := currentState.(*state.DateTimeInputState)
			if !ok {
				return nil, fmt.Errorf("invalid datetime input state: %v", currentState)
			}

			location := dateTimeInputState.Location
			if location == nil {
				location = time.Local
			}

			newState, err := convertDateTimeInputDataToState(id, &dateTimeInputData, location)
			if err != nil {
				return nil, fmt.Errorf("failed to convert date input data: %v", err)
			}

			widgetStates[id] = newState
		case state.WidgetTypeTimeInput:
			var timeInputData websocket.TimeInputData
			if err := json.Unmarshal(s, &timeInputData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal time input state: %v", err)
			}
			timeInputState, ok := currentState.(*state.TimeInputState)
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
		case state.WidgetTypeForm:
			var formData websocket.FormData
			if err := json.Unmarshal(s, &formData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal form state: %v", err)
			}
			widgetStates[id] = convertFormDataToState(&formData)
		case state.WidgetTypeButton:
			var buttonData websocket.ButtonData
			if err := json.Unmarshal(s, &buttonData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal button state: %v", err)
			}
			widgetStates[id] = convertButtonDataToState(&buttonData)
		case state.WidgetTypeMarkdown:
			var markdownData websocket.MarkdownData
			if err := json.Unmarshal(s, &markdownData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal markdown state: %v", err)
			}
			widgetStates[id] = convertMarkdownDataToState(&markdownData)
		case state.WidgetTypeColumns:
			var columnsData websocket.ColumnsData
			if err := json.Unmarshal(s, &columnsData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal columns state: %v", err)
			}
			widgetStates[id] = convertColumnsDataToState(&columnsData)
		case state.WidgetTypeColumnItem:
			var columnItemData websocket.ColumnItemData
			if err := json.Unmarshal(s, &columnItemData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal column item state: %v", err)
			}
			widgetStates[id] = convertColumnItemDataToState(&columnItemData)
		case state.WidgetTypeTable:
			var tableData websocket.TableData
			if err := json.Unmarshal(s, &tableData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal table state: %v", err)
			}
			widgetStates[id] = convertTableDataToState(id, &tableData)
		case state.WidgetTypeSelectbox:
			var selectboxData websocket.SelectboxData
			if err := json.Unmarshal(s, &selectboxData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal selectbox state: %v", err)
			}
			widgetStates[id] = convertSelectboxDataToState(id, &selectboxData)
		case state.WidgetTypeMultiSelect:
			var multiSelectData websocket.MultiSelectData
			if err := json.Unmarshal(s, &multiSelectData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal multiSelect state: %v", err)
			}
			widgetStates[id] = convertMultiSelectDataToState(id, &multiSelectData)
		case state.WidgetTypeCheckbox:
			var checkboxData websocket.CheckboxData
			if err := json.Unmarshal(s, &checkboxData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal checkbox state: %v", err)
			}
			widgetStates[id] = convertCheckboxDataToState(id, &checkboxData)
		case state.WidgetTypeCheckboxGroup:
			var checkboxGroupData websocket.CheckboxGroupData
			if err := json.Unmarshal(s, &checkboxGroupData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal checkbox group state: %v", err)
			}
			widgetStates[id] = convertCheckboxGroupDataToState(id, &checkboxGroupData)
		case state.WidgetTypeRadio:
			var radioData websocket.RadioData
			if err := json.Unmarshal(s, &radioData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal radio group state: %v", err)
			}
			widgetStates[id] = convertRadioDataToState(id, &radioData)
		case state.WidgetTypeTextArea:
			var textareaData websocket.TextAreaData
			if err := json.Unmarshal(s, &textareaData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal textarea state: %v", err)
			}
			widgetStates[id] = convertTextAreaDataToState(id, &textareaData)
		default:
			return nil, fmt.Errorf("unknown widget type: %s", currentState.GetType())
		}
	}
	return widgetStates, nil
}
