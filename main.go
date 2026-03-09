package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present; real environment variables take precedence.
	_ = godotenv.Load()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	store, err := NewPostgresStore(connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Migrate(); err != nil {
		log.Fatal("migration failed:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:9000"
	}

	h := &Handler{store: store, baseURL: baseURL}

	mux := http.NewServeMux()

	// Ticket routes
	mux.HandleFunc("POST /qrcodes", h.Create)
	mux.HandleFunc("GET /qrcodes", h.List)
	mux.HandleFunc("GET /qrcodes/phone/{phone}", h.GetByClientNumber)
	mux.HandleFunc("GET /qrcodes/{id}", h.GetByID)
	mux.HandleFunc("PATCH /qrcodes/{id}/use", h.MarkAsUsed)
	mux.HandleFunc("GET /scan/{id}", h.MarkAsUsed)
	mux.HandleFunc("GET /image/{id}", h.GetImage)

	// API documentation
	mux.HandleFunc("GET /docs/", swaggerUI)
	mux.HandleFunc("GET /docs/openapi.yaml", swaggerUI)

	addr := ":" + port
	fmt.Printf("Listening on %s\n", addr)
	fmt.Printf("API docs: http://localhost%s/docs/\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server error:", err)
	}
}
