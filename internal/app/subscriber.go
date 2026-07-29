package app

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// YieldRuntime handles GET /channels/{channel}/events.
//
// The first request for a channel creates the runtime actor: it
// registers itself in Redis, opens an SSE stream, subscribes to the
// channel, and dispatches incoming commands through the callback router.
//
// A subsequent request for the same channel sees the existing runtime
// registration, does not subscribe, and instead returns the current
// runtime info -- it never becomes a second subscriber.
func YieldRuntime(w http.ResponseWriter, r *http.Request) {
	channel := chi.URLParam(r, "channel")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE unsupported", http.StatusInternalServerError)
		return
	}

	existing, err := RedisClient.HGetAll(r.Context(), runtimeKey(channel)).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(existing) > 0 {
		// A runtime is already registered for this channel: don't
		// subscribe again, just report what's already running.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existing)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	runtime := NewRuntime(channel, w, r)

	if err := runtime.Register(); err != nil {
		return
	}

	runtime.Heartbeat()

	defer func() {
		runtime.Remove()
		runtime.Cancel()
	}()

	pubsub := RedisClient.Subscribe(runtime.Context, channel)
	defer pubsub.Close()

	for {
		select {
		case <-runtime.Context.Done():
			return

		case <-r.Context().Done():
			return

		case message, ok := <-pubsub.Channel():
			if !ok {
				return
			}

			var command Message
			if err := json.Unmarshal([]byte(message.Payload), &command); err != nil {
				continue
			}

			callback, ok := Callbacks[command.Action]
			if !ok {
				continue
			}

			callback(runtime, command)
		}
	}
}
