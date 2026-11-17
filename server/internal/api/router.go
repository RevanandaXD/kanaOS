package api

import "net/http"

func NewRouter(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", handler.Health)
	mux.HandleFunc("GET /api/services", handler.GetServices)
	mux.HandleFunc("GET /", handler.Root)

	return mux
}
