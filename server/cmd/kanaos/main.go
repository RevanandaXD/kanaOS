package main

import (
	"log" 
	"net/http"

	"github.com/revanandaxd/kanaos/server/internal/api"
  "github.com/revanandaxd/kanaos/server/internal/config"
)

func main() {
	cfg, err := config.Load("../config/service.yml")
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    handler := &api.Handler{
        Config: cfg,
    }

    router := api.NewRouter(handler)

    log.Println("KanaOS API running on :8080")

    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal("Server error:", err)
    }
}
