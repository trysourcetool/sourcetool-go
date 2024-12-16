package ui

import (
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/runtime"
	"github.com/trysourcetool/sourcetool-go/ui/textinput"
	ws "github.com/trysourcetool/sourcetool-go/websocket"
)

func TextInput(label string, options ...textinput.Option) string {
	opts := &textinput.Options{
		Label: label,
	}

	for _, option := range options {
		option(opts)
	}

	var returnValue string
	// TODO: check state and assign to returnValue

	// queue message to runtime
	if runtime.Runtime != nil {
		runtime.Runtime.EnqueueMessage(uuid.Must(uuid.NewV4()).String(), ws.MessageMethodInitializeHost, opts)
	}

	return returnValue
}
