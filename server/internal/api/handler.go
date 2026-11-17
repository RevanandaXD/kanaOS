package api

import (
	"encoding/json"
	"net/http"

	"github.com/revanandaxd/kanaos/server/internal/config"
)

type Handler struct {
	Config *config.Config
}

func (h *Handler) Root(w http.ResponseWriter, r * http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"title": "Welcome to Endpoint API",
	})
}

func (h *Handler) GetServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(h.Config)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"name":   "KanaOS API",
	})
}
