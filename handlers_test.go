package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerCreate(t *testing.T) {
	store := NewMemStore()
	h := &Handler{store: store, baseURL: "http://localhost:9000"}

	// First test - success case
	req := httptest.NewRequest(http.MethodPost, "/qrcodes", strings.NewReader(`{"client_number":"123456"}`))
	w := httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Create() status = %v, want %v", w.Code, http.StatusCreated)
	}

	// Test invalid body
	req = httptest.NewRequest(http.MethodPost, "/qrcodes", strings.NewReader(`not json`))
	w = httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	// Test empty client_number
	req = httptest.NewRequest(http.MethodPost, "/qrcodes", strings.NewReader(`{"client_number":""}`))
	w = httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	// Test non-numeric client_number
	req = httptest.NewRequest(http.MethodPost, "/qrcodes", strings.NewReader(`{"client_number":"abc"}`))
	w = httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Create() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	// Test duplicate client_number - this should fail because we already created "123456"
	req = httptest.NewRequest(http.MethodPost, "/qrcodes", strings.NewReader(`{"client_number":"123456"}`))
	w = httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("Create() status = %v, want %v, body: %s", w.Code, http.StatusConflict, w.Body.String())
	}
}

func TestHandlerList(t *testing.T) {
	store := NewMemStore()
	store.Create(QRCode{ID: "1", ClientNumber: "111", Used: false})
	store.Create(QRCode{ID: "2", ClientNumber: "222", Used: false})

	h := &Handler{store: store, baseURL: "http://localhost:9000"}

	req := httptest.NewRequest(http.MethodGet, "/qrcodes", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("List() status = %v, want %v", w.Code, http.StatusOK)
	}

	var list []QRCode
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("List() returned %d items, want 2", len(list))
	}
}

func TestHandlerGetByID(t *testing.T) {
	store := NewMemStore()
	store.Create(QRCode{ID: "test-id", ClientNumber: "123456", Used: false})

	h := &Handler{store: store, baseURL: "http://localhost:9000"}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /qrcodes/{id}", h.GetByID)

	// Test found
	req := httptest.NewRequest(http.MethodGet, "/qrcodes/test-id", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GetByID() status = %v, want %v, body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	// Test not found
	req = httptest.NewRequest(http.MethodGet, "/qrcodes/nonexistent", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("GetByID() status = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestHandlerMarkAsUsed(t *testing.T) {
	store := NewMemStore()
	store.Create(QRCode{ID: "test-id", ClientNumber: "123456", Used: false})
	store.Create(QRCode{ID: "used-id", ClientNumber: "999999", Used: true})

	h := &Handler{store: store, baseURL: "http://localhost:9000"}

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /qrcodes/{id}/use", h.MarkAsUsed)

	// Test success
	req := httptest.NewRequest(http.MethodPatch, "/qrcodes/test-id/use", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("MarkAsUsed() status = %v, want %v, body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	// Test not found
	req = httptest.NewRequest(http.MethodPatch, "/qrcodes/nonexistent/use", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("MarkAsUsed() status = %v, want %v", w.Code, http.StatusNotFound)
	}

	// Test already used
	req = httptest.NewRequest(http.MethodPatch, "/qrcodes/used-id/use", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("MarkAsUsed() status = %v, want %v", w.Code, http.StatusConflict)
	}
}

func TestHandlerDelete(t *testing.T) {
	store := NewMemStore()
	store.Create(QRCode{ID: "test-id", ClientNumber: "123456", Used: false})

	h := &Handler{store: store, baseURL: "http://localhost:9000"}

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /qrcodes/{id}", h.Delete)

	// Test success
	req := httptest.NewRequest(http.MethodDelete, "/qrcodes/test-id", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("Delete() status = %v, want %v", w.Code, http.StatusNoContent)
	}

	// Test not found
	req = httptest.NewRequest(http.MethodDelete, "/qrcodes/nonexistent", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Delete() status = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestHandlerScan(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		wantStatus   int
		wantContains []string
		setup        func(*MemStore)
	}{
		{
			name:         "valid ticket",
			id:           "test-id",
			wantStatus:   http.StatusOK,
			wantContains: []string{"Valid Ticket", "Client #: 123456"},
			setup: func(s *MemStore) {
				s.Create(QRCode{ID: "test-id", ClientNumber: "123456", Used: false})
			},
		},
		{
			name:         "not found",
			id:           "nonexistent",
			wantStatus:   http.StatusNotFound,
			wantContains: []string{"Ticket Not Found"},
			setup:        func(s *MemStore) {},
		},
		{
			name:         "already used",
			id:           "used-id",
			wantStatus:   http.StatusConflict,
			wantContains: []string{"Already Used"},
			setup: func(s *MemStore) {
				s.Create(QRCode{ID: "used-id", ClientNumber: "999999", Used: true})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemStore()
			tt.setup(store)

			h := &Handler{store: store, baseURL: "http://localhost:9000"}

			mux := http.NewServeMux()
			mux.HandleFunc("GET /scan/{id}", h.Scan)

			req := httptest.NewRequest(http.MethodGet, "/scan/"+tt.id, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Scan() status = %v, want %v, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}

			for _, s := range tt.wantContains {
				if !strings.Contains(w.Body.String(), s) {
					t.Errorf("Scan() body should contain %q, got: %s", s, w.Body.String())
				}
			}
		})
	}
}

func TestHandlerGetImage(t *testing.T) {
	store := NewMemStore()
	qr := QRCode{ID: "test-id", ClientNumber: "123456", Used: false, Image: "http://localhost:9000/scan/test-id"}
	store.Create(qr)

	h := &Handler{store: store, baseURL: "http://localhost:9000"}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /image/{id}", h.GetImage)

	req := httptest.NewRequest(http.MethodGet, "/image/test-id", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetImage() status = %v, want %v. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	if w.Header().Get("Content-Type") != "image/png" {
		t.Errorf("GetImage() Content-Type = %v, want image/png", w.Header().Get("Content-Type"))
	}

	if w.Header().Get("Content-Disposition") != "attachment; filename=\"qrcode.png\"" {
		t.Errorf("GetImage() Content-Disposition = %v, want attachment; filename=\"qrcode.png\"", w.Header().Get("Content-Disposition"))
	}
}

func TestHandlerGetImageNotFound(t *testing.T) {
	store := NewMemStore()
	h := &Handler{store: store, baseURL: "http://localhost:9000"}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /image/{id}", h.GetImage)

	req := httptest.NewRequest(http.MethodGet, "/image/nonexistent", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GetImage() status = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestHandlerCreateResponse(t *testing.T) {
	store := NewMemStore()
	h := &Handler{store: store, baseURL: "http://localhost:9000"}

	body := bytes.NewBufferString(`{"client_number":"123456"}`)
	req := httptest.NewRequest(http.MethodPost, "/qrcodes", body)
	w := httptest.NewRecorder()

	h.Create(w, req)

	var qr QRCode
	if err := json.Unmarshal(w.Body.Bytes(), &qr); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if qr.ID == "" {
		t.Error("ID should be auto-generated")
	}
	if qr.ClientNumber != "123456" {
		t.Errorf("ClientNumber = %v, want 123456", qr.ClientNumber)
	}
	if qr.Image == "" {
		t.Error("Image should be set")
	}
	if qr.Used != false {
		t.Errorf("Used = %v, want false", qr.Used)
	}
	if qr.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}
