package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	store, err := NewPostgresStore(connStr)
	if err != nil {
		log.Fatal(err)
	}

	h := &Handler{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /qrcodes", h.Create)
	mux.HandleFunc("GET /qrcodes", h.List)
	mux.HandleFunc("GET /qrcodes/phone/{phone}", h.GetByClientNumber)
	mux.HandleFunc("GET /qrcodes/{id}", h.GetByID)

	fmt.Println("Listening on :9000")
	if err := http.ListenAndServe(":9000", mux); err != nil {
		fmt.Println("Server error:", err)
	}
}
