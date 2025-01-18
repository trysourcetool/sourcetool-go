package websocket

import (
	"encoding/json"
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
	MessageMethodRerunPage        MessageMethod = "RERUN_PAGE"
	MessageMethodCloseSession     MessageMethod = "CLOSE_SESSION"
)

type MessageHandlerFunc func(*Message) error

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
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Route  string   `json:"route"`
	Path   []int    `json:"path"`
	Groups []string `json:"groups"`
}

type RenderWidgetPayload struct {
	SessionID  string `json:"sessionId"`
	PageID     string `json:"pageId"`
	WidgetID   string `json:"widgetId"`
	WidgetType string `json:"widgetType"`
	Path       []int  `json:"path"`
	Data       any    `json:"data"`
}

type InitializeClientPayload struct {
	SessionID string `json:"sessionId"`
	PageID    string `json:"pageId"`
}

type RerunPagePayload struct {
	SessionID string          `json:"sessionId"`
	PageID    string          `json:"pageId"`
	State     json.RawMessage `json:"state"`
}

type CloseSessionPayload struct {
	SessionID string `json:"sessionId"`
}

type TextInputData struct {
	Value        string `json:"value"`
	Label        string `json:"label"`
	Placeholder  string `json:"placeholder"`
	DefaultValue string `json:"defaultValue"`
	Required     bool   `json:"required"`
	Disabled     bool   `json:"disabled"`
	MaxLength    *int   `json:"maxLength"`
	MinLength    *int   `json:"minLength"`
}

type NumberInputData struct {
	Value        *float64 `json:"value"`
	Label        string   `json:"label"`
	Placeholder  string   `json:"placeholder"`
	DefaultValue *float64 `json:"defaultValue"`
	Required     bool     `json:"required"`
	Disabled     bool     `json:"disabled"`
	MaxValue     *float64 `json:"maxValue"`
	MinValue     *float64 `json:"minValue"`
}

type DateInputData struct {
	Value        string `json:"value"`
	Label        string `json:"label"`
	Placeholder  string `json:"placeholder"`
	DefaultValue string `json:"defaultValue"`
	Required     bool   `json:"required"`
	Disabled     bool   `json:"disabled"`
	Format       string `json:"format"`
	MaxValue     string `json:"maxValue"`
	MinValue     string `json:"minValue"`
}

type DateTimeInputData struct {
	Value        string `json:"value"`
	Label        string `json:"label"`
	Placeholder  string `json:"placeholder"`
	DefaultValue string `json:"defaultValue"`
	Required     bool   `json:"required"`
	Disabled     bool   `json:"disabled"`
	Format       string `json:"format"`
	MaxValue     string `json:"maxValue"`
	MinValue     string `json:"minValue"`
}

type TimeInputData struct {
	Value        string `json:"value"`
	Label        string `json:"label"`
	Placeholder  string `json:"placeholder"`
	DefaultValue string `json:"defaultValue"`
	Required     bool   `json:"required"`
	Disabled     bool   `json:"disabled"`
}

type SelectboxData struct {
	Value        *int     `json:"value"`
	Label        string   `json:"label"`
	Options      []string `json:"options"`
	Placeholder  string   `json:"placeholder"`
	DefaultValue *int     `json:"defaultValue"`
	Required     bool     `json:"required"`
	Disabled     bool     `json:"disabled"`
}

type MultiSelectData struct {
	Value        []int    `json:"value"`
	Label        string   `json:"label"`
	Options      []string `json:"options"`
	Placeholder  string   `json:"placeholder"`
	DefaultValue []int    `json:"defaultValue"`
	Required     bool     `json:"required"`
	Disabled     bool     `json:"disabled"`
}

type CheckboxData struct {
	Value        bool   `json:"value"`
	Label        string `json:"label"`
	DefaultValue bool   `json:"defaultValue"`
	Required     bool   `json:"required"`
	Disabled     bool   `json:"disabled"`
}

type CheckboxGroupData struct {
	Value        []int    `json:"value"`
	Label        string   `json:"label"`
	Options      []string `json:"options"`
	DefaultValue []int    `json:"defaultValue"`
	Required     bool     `json:"required"`
	Disabled     bool     `json:"disabled"`
}

type RadioData struct {
	Value        *int     `json:"value"`
	Label        string   `json:"label"`
	Options      []string `json:"options"`
	DefaultValue *int     `json:"defaultValue"`
	Required     bool     `json:"required"`
	Disabled     bool     `json:"disabled"`
}

type TextAreaData struct {
	Value        string `json:"value"`
	Label        string `json:"label"`
	Placeholder  string `json:"placeholder"`
	DefaultValue string `json:"defaultValue"`
	Required     bool   `json:"required"`
	Disabled     bool   `json:"disabled"`
	MaxLength    *int   `json:"maxLength"`
	MinLength    *int   `json:"minLength"`
	MaxLines     *int   `json:"maxLines"`
	MinLines     *int   `json:"minLines"`
	AutoResize   bool   `json:"autoResize"`
}

type FormData struct {
	Value          bool   `json:"value"`
	ButtonLabel    string `json:"buttonLabel"`
	ButtonDisabled bool   `json:"buttonDisabled"`
	ClearOnSubmit  bool   `json:"clearOnSubmit"`
}

type ButtonData struct {
	Value    bool   `json:"value"`
	Label    string `json:"label"`
	Disabled bool   `json:"disabled"`
}

type MarkdownData struct {
	Body string `json:"body"`
}

type TableData struct {
	Data         any            `json:"data"`
	Value        TableDataValue `json:"value"`
	Header       *string        `json:"header"`
	Description  *string        `json:"description"`
	OnSelect     *string        `json:"onSelect"`
	RowSelection *string        `json:"rowSelection"`
}

type TableDataValue struct {
	Selection *TableDataValueSelection `json:"selection"`
}

type TableDataValueSelection struct {
	Row  int   `json:"row"`
	Rows []int `json:"rows"`
}

type ColumnsData struct {
	Columns int `json:"columns"`
}

type ColumnItemData struct {
	Weight float64 `json:"weight"`
}
