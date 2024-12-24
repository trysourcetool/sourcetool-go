package sourcetool

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofrs/uuid/v5"

	ws "github.com/trysourcetool/sourcetool-go/websocket"
)

// initializeClientHandler handles INITIALIZE_CLIENT messages
type initializeClientHandler struct {
	r *runtime
}

func (h *initializeClientHandler) Handle(msg *ws.Message) error {
	var p ws.InitializeClientPayload
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
	session := NewSession(sessionID)
	h.r.sessionManager.SetSession(session)

	page := h.r.pageManager.GetPage(pageID)
	if page == nil {
		return fmt.Errorf("page not found: %s", pageID)
	}

	ui := &uiBuilder{
		context: context.Background(),
		runtime: h.r,
		session: session,
		page:    page,
		cursor:  NewCursor(MAIN),
	}

	if err := page.Run(ui); err != nil {
		return fmt.Errorf("failed to run page: %v", err)
	}

	return nil
}

type rerunPageHandler struct {
	r *runtime
}

func (h *rerunPageHandler) Handle(msg *ws.Message) error {
	var p ws.RerunPagePayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	sessionID, err := uuid.FromString(p.SessionID)
	if err != nil {
		return err
	}
	sess := h.r.sessionManager.GetSession(sessionID)
	if sess == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	pageID, err := uuid.FromString(p.PageID)
	if err != nil {
		return err
	}
	page := h.r.pageManager.GetPage(pageID)
	if page == nil {
		return fmt.Errorf("page not found: %s", pageID)
	}

	var states map[uuid.UUID]any
	if err := json.Unmarshal(p.State, &states); err != nil {
		return fmt.Errorf("failed to unmarshal state: %v", err)
	}
	sess.State.SetStates(states)

	ui := &uiBuilder{
		context: context.Background(),
		runtime: h.r,
		session: sess,
		page:    page,
		cursor:  NewCursor(MAIN),
	}

	if err := page.Run(ui); err != nil {
		return fmt.Errorf("failed to run page: %v", err)
	}

	sess.State.ResetButtons()

	return nil
}

// closeSessionHandler handles CLOSE_SESSION messages
type closeSessionHandler struct {
	r *runtime
}

func (h *closeSessionHandler) Handle(msg *ws.Message) error {
	var p ws.CloseSessionPayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	sessionID, err := uuid.FromString(p.SessionID)
	if err != nil {
		return err
	}

	h.r.sessionManager.DeleteSession(sessionID)

	return nil
}
