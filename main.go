package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	_ = godotenv.Load()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	store, err := NewPostgresStore(connStr)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	if err := store.Migrate(); err != nil {
		slog.Error("database migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("database connected and migrated")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}

	h := &Handler{store: store, baseURL: baseURL}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /qrcodes", h.Create)
	mux.HandleFunc("GET /qrcodes", h.List)
	mux.HandleFunc("GET /qrcodes/phone/{phone}", h.GetByClientNumber)
	mux.HandleFunc("GET /qrcodes/{id}", h.GetByID)
	mux.HandleFunc("PATCH /qrcodes/{id}/use", h.MarkAsUsed)
	mux.HandleFunc("DELETE /qrcodes/{id}", h.Delete)
	mux.HandleFunc("GET /scan/{id}", h.Scan)
	mux.HandleFunc("GET /image/{id}", h.GetImage)

	mux.HandleFunc("GET /docs/", swaggerUI)
	mux.HandleFunc("GET /docs/openapi.yaml", swaggerUI)

	addr := ":" + port

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		slog.Info("starting server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("server started", "docs_url", "http://localhost"+addr+"/docs/")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}
