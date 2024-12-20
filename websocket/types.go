package websocket

import (
	"encoding/json"
	"errors"
)

type MessageKind string

const (
	MessageKindCall     MessageKind = "CALL"
	MessageKindResponse MessageKind = "RESPONSE"
)

type MessageMethod string

const (
	MessageMethodInitializeHost   MessageMethod = "INITIALIZE_HOST"
	MessageMethodInitializeClient MessageMethod = "INITIALIZE_CLIENT"
	MessageMethodRenderWidget     MessageMethod = "RENDER_WIDGET"
	MessageMethodCloseSession     MessageMethod = "CLOSE_SESSION"
)

// MessageHandler handles websocket messages
type MessageHandler interface {
	Handle(*Message) error
}

// Common errors
var (
	ErrUnknownMethod = errors.New("unknown message method")
)

type Message struct {
	ID      string          `json:"id"`
	Kind    MessageKind     `json:"kind"`
	Method  MessageMethod   `json:"method"`
	Payload json.RawMessage `json:"payload"`
	Error   *MessageError   `json:"error"`
}

type MessageError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type InitializeHostPayload struct {
	APIKey     string                       `json:"apiKey"`
	SDKName    string                       `json:"sdkName"`
	SDKVersion string                       `json:"sdkVersion"`
	Pages      []*InitializeHostPagePayload `json:"pages"`
}

type InitializeHostPagePayload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RenderWidgetPayload struct {
	SessionID  string `json:"sessionId"`
	PageID     string `json:"pageId"`
	WidgetID   string `json:"widgetId"`
	WidgetType string `json:"widgetType"`
	Data       any    `json:"data"`
}

type InitializeClientPayload struct {
	SessionID string `json:"sessionId"`
	PageID    string `json:"pageId"`
}

type CloseSessionPayload struct {
	SessionID string `json:"sessionId"`
}
