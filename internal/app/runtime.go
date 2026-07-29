package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Runtime is the Cloud Function runtime actor. One Runtime exists per
// live invocation of the /events endpoint for a given channel. It owns
// cancellation, the Redis registration/subscription, the SSE writer, and
// cleanup on exit.
type Runtime struct {
	ID      string
	Channel string

	Context context.Context
	Cancel  context.CancelFunc

	Writer  http.ResponseWriter
	Flusher http.Flusher
}

func runtimeKey(channel string) string {
	return "runtime:" + channel
}

func NewRuntime(channel string, w http.ResponseWriter, r *http.Request) *Runtime {
	ctx, cancel := context.WithCancel(r.Context())

	id := os.Getenv("K_REVISION")
	if id == "" {
		id = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	flusher := w.(http.Flusher)

	return &Runtime{
		ID:      id,
		Channel: channel,
		Context: ctx,
		Cancel:  cancel,
		Writer:  w,
		Flusher: flusher,
	}
}

// Register writes this runtime's registration into Redis, marking it as
// the active subscriber for its channel.
func (rt *Runtime) Register() error {
	return RedisClient.HSet(
		rt.Context,
		runtimeKey(rt.Channel),
		map[string]interface{}{
			"id":      rt.ID,
			"status":  "ACTIVE",
			"started": time.Now().Unix(),
		},
	).Err()
}

// Remove clears the runtime's Redis registration. Uses a background
// context since it is typically called during cleanup after rt.Context
// has already been cancelled.
func (rt *Runtime) Remove() {
	RedisClient.Del(
		context.Background(),
		runtimeKey(rt.Channel),
	)
}

// Heartbeat periodically refreshes the registration's TTL so Redis stays
// aware that the runtime is still alive. It stops when rt.Context is
// cancelled.
func (rt *Runtime) Heartbeat() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-rt.Context.Done():
				return

			case <-ticker.C:
				RedisClient.Expire(
					context.Background(),
					runtimeKey(rt.Channel),
					90*time.Second,
				)
			}
		}
	}()
}

// Send writes a single SSE event to the runtime's HTTP response and
// flushes it immediately.
func (rt *Runtime) Send(event string, data string) {
	fmt.Fprintf(rt.Writer, "event: %s\n", event)
	fmt.Fprintf(rt.Writer, "data: %s\n\n", data)
	rt.Flusher.Flush()
}
