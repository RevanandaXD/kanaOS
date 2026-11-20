package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/revanandaxd/kanaos/server/internal/api"
	"github.com/revanandaxd/kanaos/server/internal/config"
	"github.com/revanandaxd/kanaos/server/internal/docker"
)

func main() {
	configPath := env("CONFIG_PATH", "../config/service.yml")
	cfg, err := config.Load(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("Config %q tidak ditemukan, melanjutkan tanpa service config", configPath)
			cfg = &config.Config{}
		} else {
			log.Fatalf("Failed to load config: %v", err)
		}
	}

	dockerClient, err := docker.New()
	if err != nil {
		log.Fatalf("Failed to create docker client: %v", err)
	}
	defer dockerClient.Close()

	pingCtx, cancelPing := context.WithTimeout(context.Background(), 3*time.Second)
	if err := dockerClient.Ping(pingCtx); err != nil {
		log.Printf("Warning: Docker is not reachable: %v", err)
	}
	cancelPing()

	handler := &api.Handler{
		Config: cfg,
		Cache:  docker.NewCache(dockerClient, duration("DOCKER_CACHE_TTL", 3*time.Second)),
	}

	srv := &http.Server{
		Addr:              ":" + env("PORT", "8080"),
		Handler:           api.NewRouter(handler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("KanaOS API running on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("KanaOS API stopped")
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func duration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("Warning: invalid %s=%q, memakai default %s", key, v, fallback)
		return fallback
	}
	return d
}
