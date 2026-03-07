package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	store Store
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var qr QRCode
	if err := json.NewDecoder(r.Body).Decode(&qr); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.store.Create(qr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(qr)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	qrs, err := h.store.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(qrs)
}

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

func (h *Handler) GetByClientNumber(w http.ResponseWriter, r *http.Request) {
	phone := r.PathValue("phone")
	qr, err := h.store.GetByClientNumber(phone)
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
