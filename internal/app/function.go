package app

import "net/http"

// Router is the HTTP surface used by both the Cloud Function entrypoint
// (Function) and any standalone server (see cmd/api).
var Router http.Handler

func init() {
	Router = NewRouter()
}

// Function is the Cloud Function entrypoint. Each invocation is handed to
// the shared router; the runtime actor lifecycle is driven from there.
func Function(w http.ResponseWriter, r *http.Request) {
	Router.ServeHTTP(w, r)
}
