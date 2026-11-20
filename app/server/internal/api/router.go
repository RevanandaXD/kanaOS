package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func NewRouter(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", handler.Health)
	mux.HandleFunc("GET /api/services", handler.GetServices)
	mux.HandleFunc("GET /api/containers", handler.GetContainers)
	mux.HandleFunc("GET /api/system", handler.GetSystem)

	upstream := os.Getenv("WEB_UPSTREAM")
	if upstream == "" {
		upstream = "http://127.0.0.1:3000"
	}
	target, err := url.Parse(upstream)
	if err != nil {
		target, _ = url.Parse("http://127.0.0.1:3000")
	}
	mux.Handle("/", httputil.NewSingleHostReverseProxy(target))

	return mux
}
