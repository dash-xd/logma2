package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/channels/{channel}/events", YieldRuntime)

	r.Post("/channels/{channel}/publish", Publish)

	return r
}
