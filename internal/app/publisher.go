package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Publish handles POST /channels/{channel}/publish. Publishers never
// subscribe -- they only send commands for the channel's runtime actor
// to dispatch through the callback router.
func Publish(w http.ResponseWriter, r *http.Request) {
	channel := chi.URLParam(r, "channel")

	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := PublishMessage(r.Context(), channel, message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func PublishMessage(ctx context.Context, channel string, message Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return RedisClient.Publish(ctx, channel, payload).Err()
}
