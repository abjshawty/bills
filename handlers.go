package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/google/uuid"
)

// Handler holds the Store and exposes HTTP handler methods.
type Handler struct {
	store   Store
	baseURL string
}

// Create handles POST /qrcodes.
// It decodes a QRCode from the request body, generates a UUID for the ID,
// derives the image URL from BASE_URL, persists it, and returns the created ticket as JSON.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var qr QRCode
	if err := json.NewDecoder(r.Body).Decode(&qr); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	qr.ID = uuid.New().String()
	qr.Image = h.baseURL + "/scan/" + qr.ID
	if err := h.store.Create(qr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(qr)
}

// List handles GET /qrcodes.
// It returns all stored QR code tickets as a JSON array.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	qrs, err := h.store.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(qrs)
}

// GetByID handles GET /qrcodes/{id}.
// It returns the ticket with the matching ID, or 404 if not found.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	qr, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(qr)
}

// GetByClientNumber handles GET /qrcodes/phone/{phone}.
// It returns all tickets associated with the given phone number, or 404 if none found.
func (h *Handler) GetByClientNumber(w http.ResponseWriter, r *http.Request) {
	phone := r.PathValue("phone")
	qrs, err := h.store.GetByClientNumber(phone)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(qrs)
}

// MarkAsUsed handles PATCH /qrcodes/{id}/use.
// It marks the ticket as used and returns the updated ticket as JSON.
func (h *Handler) MarkAsUsed(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.MarkAsUsed(id); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	qr, err := h.store.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(qr)
}
