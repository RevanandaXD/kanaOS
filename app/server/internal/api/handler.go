package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/revanandaxd/kanaos/server/internal/config"
	"github.com/revanandaxd/kanaos/server/internal/docker"
	"github.com/revanandaxd/kanaos/server/internal/system"
)

type Handler struct {
	Config *config.Config
	Cache  *docker.Cache
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetContainers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	containers, err := h.Cache.Containers(ctx)
	if err != nil {
		log.Println("docker list:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "docker tidak bisa dihubungi",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"containers": containers,
	})
}

func (h *Handler) GetSystem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(system.Collect())
}
