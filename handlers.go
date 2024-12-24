package sourcetool

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/websocket"
)

type messageHandler struct {
	r *runtime
}

func (h *messageHandler) initializeCilent(msg *websocket.Message) error {
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
	session := newSession(sessionID)
	h.r.sessionManager.setSession(session)

	page := h.r.pageManager.getPage(pageID)
	if page == nil {
		return fmt.Errorf("page not found: %s", pageID)
	}

	ui := &uiBuilder{
		context: context.Background(),
		runtime: h.r,
		session: session,
		page:    page,
		cursor:  newCursor(main),
	}

	if err := page.run(ui); err != nil {
		return fmt.Errorf("failed to run page: %v", err)
	}

	return nil
}

func (h *messageHandler) rerunPage(msg *websocket.Message) error {
	var p websocket.RerunPagePayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	sessionID, err := uuid.FromString(p.SessionID)
	if err != nil {
		return err
	}
	sess := h.r.sessionManager.getSession(sessionID)
	if sess == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	pageID, err := uuid.FromString(p.PageID)
	if err != nil {
		return err
	}
	page := h.r.pageManager.getPage(pageID)
	if page == nil {
		return fmt.Errorf("page not found: %s", pageID)
	}

	var states map[uuid.UUID]any
	if err := json.Unmarshal(p.State, &states); err != nil {
		return fmt.Errorf("failed to unmarshal state: %v", err)
	}
	sess.state.setStates(states)

	ui := &uiBuilder{
		context: context.Background(),
		runtime: h.r,
		session: sess,
		page:    page,
		cursor:  newCursor(main),
	}

	if err := page.run(ui); err != nil {
		return fmt.Errorf("failed to run page: %v", err)
	}

	sess.state.resetButtons()

	return nil
}

func (h *messageHandler) closeSession(msg *websocket.Message) error {
	var p websocket.CloseSessionPayload
	if err := json.Unmarshal(msg.Payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	sessionID, err := uuid.FromString(p.SessionID)
	if err != nil {
		return err
	}

	h.r.sessionManager.deleteSession(sessionID)

	return nil
}
