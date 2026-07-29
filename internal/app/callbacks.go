package app

import "encoding/json"

// Message is the wire format published to a channel and dispatched to a
// callback. Commands are not necessarily HTTP responses, so callbacks
// receive the Runtime rather than an http.ResponseWriter.
type Message struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

type Callback func(*Runtime, Message)

var Callbacks = map[string]Callback{
	"echo":     Echo,
	"shutdown": Shutdown,
}

func Echo(rt *Runtime, msg Message) {
	data, _ := json.Marshal(msg.Data)
	rt.Send("message", string(data))
}

// Shutdown gracefully terminates the invocation: it notifies the SSE
// client, then cancels the runtime context, which unwinds the subscriber
// loop and triggers cleanup.
func Shutdown(rt *Runtime, msg Message) {
	rt.Send("shutdown", "runtime stopping")
	rt.Cancel()
}
