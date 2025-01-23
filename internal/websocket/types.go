package websocket

import (
	"fmt"

	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type MessageHandlerFunc func(*websocketv1.Message) error

func UnmarshalMessage(data []byte) (*websocketv1.Message, error) {
	var msg websocketv1.Message
	if err := protojson.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func MarshalMessage(msg *websocketv1.Message) ([]byte, error) {
	return protojson.Marshal(msg)
}

func NewMessage(id string, payload proto.Message) (*websocketv1.Message, error) {
	msg := &websocketv1.Message{
		Id: id,
	}

	switch p := payload.(type) {
	case *websocketv1.InitializeHost:
		msg.Type = &websocketv1.Message_InitializeHost{InitializeHost: p}
	case *websocketv1.InitializeClient:
		msg.Type = &websocketv1.Message_InitializeClient{InitializeClient: p}
	case *websocketv1.RenderWidget:
		msg.Type = &websocketv1.Message_RenderWidget{RenderWidget: p}
	case *websocketv1.RerunPage:
		msg.Type = &websocketv1.Message_RerunPage{RerunPage: p}
	case *websocketv1.CloseSession:
		msg.Type = &websocketv1.Message_CloseSession{CloseSession: p}
	default:
		return nil, fmt.Errorf("unsupported message type: %T", payload)
	}

	return msg, nil
}
