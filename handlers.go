package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	qrcodegen "github.com/skip2/go-qrcode"
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

	if qr.ClientNumber == "" || !isValidClientNumber(qr.ClientNumber) {
		http.Error(w, "client_number must be a non-empty numeric string", http.StatusBadRequest)
		return
	}

	qr.ID = uuid.New().String()
	qr.CreatedAt = time.Now()
	qr.Image = h.baseURL + "/image/" + qr.ID
	if err := h.store.Create(qr); err != nil {
		if errors.Is(err, ErrDuplicateClientNumber) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(qr)
}

func isValidClientNumber(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.Atoi(s)
	return err == nil
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

// GetImage handles GET /image/{id}.
// It returns a PNG of the QR code whose content is the base64-encoded scan URL.
func (h *Handler) GetImage(w http.ResponseWriter, r *http.Request) {
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

	png, err := qrcodegen.Encode(qr.Image, qrcodegen.Medium, 256)
	if err != nil {
		http.Error(w, "failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=\"qrcode.png\"")
	w.Write(png)
}

// Scan handles GET /scan/{id}.
// It marks the ticket as used and returns a styled HTML page showing the scan result.
func (h *Handler) Scan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	qr, err := h.store.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(h.scanPage("Ticket Not Found", "The ticket ID is invalid or does not exist.", false)))
		return
	}

	if err := h.store.MarkAsUsed(id); err != nil {
		if errors.Is(err, ErrAlreadyUsed) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(h.scanPage("Already Used", "This ticket has already been scanned.", false)))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(h.scanPage("Valid Ticket", "Client #: "+qr.ClientNumber, true)))
}

func (h *Handler) scanPage(title, message string, success bool) string {
	color := "#22c55e"
	icon := "✓"
	if !success {
		color = "#ef4444"
		icon = "✗"
	}
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>` + title + `</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; background: #fafafa; }
    .card { background: white; padding: 2rem; border-radius: 1rem; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; }
    .status { width: 80px; height: 80px; border-radius: 50%; background: ` + color + `; margin: 0 auto 1rem; display: flex; align-items: center; justify-content: center; }
    .status-icon { color: white; font-size: 2.5rem; }
    h1 { margin: 0 0 0.5rem; color: #1f2937; }
    p { margin: 0; color: #6b7280; }
  </style>
</head>
<body>
  <div class="card">
    <div class="status"><span class="status-icon">` + icon + `</span></div>
    <h1>` + title + `</h1>
    <p>` + message + `</p>
  </div>
</body>
</html>`
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
		if errors.Is(err, ErrAlreadyUsed) {
			http.Error(w, err.Error(), http.StatusConflict)
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

// Delete handles DELETE /qrcodes/{id}.
// It removes the ticket from the store.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
