package main

import (
	"testing"
	"time"
)

func TestMemStoreCreate(t *testing.T) {
	tests := []struct {
		name    string
		qr      QRCode
		wantErr error
	}{
		{
			name: "success",
			qr: QRCode{
				ID:           "test-id",
				ClientNumber: "123456",
				Used:         false,
			},
			wantErr: nil,
		},
		{
			name: "duplicate client number",
			qr: QRCode{
				ID:           "test-id-2",
				ClientNumber: "123456",
				Used:         false,
			},
			wantErr: ErrDuplicateClientNumber,
		},
	}

	store := NewMemStore()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Create(tt.qr)
			if err != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStoreList(t *testing.T) {
	store := NewMemStore()

	store.Create(QRCode{ID: "1", ClientNumber: "111", Used: false})
	store.Create(QRCode{ID: "2", ClientNumber: "222", Used: false})

	list, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(list) != 2 {
		t.Errorf("List() returned %d items, want 2", len(list))
	}
}

func TestMemStoreGetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "found",
			id:      "test-id",
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      "nonexistent",
			wantErr: ErrNotFound,
		},
	}

	store := NewMemStore()
	store.Create(QRCode{ID: "test-id", ClientNumber: "123456", Used: false})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := store.GetByID(tt.id)
			if err != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStoreGetByClientNumber(t *testing.T) {
	tests := []struct {
		name      string
		phone     string
		wantCount int
		wantErr   error
		setup     func(*MemStore)
	}{
		{
			name:      "found",
			phone:     "123456",
			wantCount: 1,
			wantErr:   nil,
			setup: func(s *MemStore) {
				s.Create(QRCode{ID: "1", ClientNumber: "123456", Used: false})
			},
		},
		{
			name:      "not found",
			phone:     "999999",
			wantCount: 0,
			wantErr:   ErrNotFound,
			setup:     func(s *MemStore) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemStore()
			tt.setup(store)

			list, err := store.GetByClientNumber(tt.phone)
			if err != tt.wantErr {
				t.Errorf("GetByClientNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(list) != tt.wantCount {
				t.Errorf("GetByClientNumber() returned %d items, want %d", len(list), tt.wantCount)
			}
		})
	}
}

func TestMemStoreMarkAsUsed(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr error
		setup   func(*MemStore)
	}{
		{
			name:    "success",
			id:      "test-id",
			wantErr: nil,
			setup: func(s *MemStore) {
				s.Create(QRCode{ID: "test-id", ClientNumber: "123456", Used: false})
			},
		},
		{
			name:    "not found",
			id:      "nonexistent",
			wantErr: ErrNotFound,
			setup:   func(s *MemStore) {},
		},
		{
			name:    "already used",
			id:      "used-id",
			wantErr: ErrAlreadyUsed,
			setup: func(s *MemStore) {
				s.Create(QRCode{ID: "used-id", ClientNumber: "999999", Used: true})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemStore()
			tt.setup(store)

			err := store.MarkAsUsed(tt.id)
			if err != tt.wantErr {
				t.Errorf("MarkAsUsed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStoreDelete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr error
		setup   func(*MemStore)
	}{
		{
			name:    "success",
			id:      "test-id",
			wantErr: nil,
			setup: func(s *MemStore) {
				s.Create(QRCode{ID: "test-id", ClientNumber: "123456", Used: false})
			},
		},
		{
			name:    "not found",
			id:      "nonexistent",
			wantErr: ErrNotFound,
			setup:   func(s *MemStore) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemStore()
			tt.setup(store)

			err := store.Delete(tt.id)
			if err != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				_, err := store.GetByID(tt.id)
				if err != ErrNotFound {
					t.Errorf("Delete() did not remove the ticket")
				}
			}
		})
	}
}

func TestMemStoreTimestamps(t *testing.T) {
	store := NewMemStore()
	before := time.Now()

	store.Create(QRCode{ID: "test-id", ClientNumber: "123456", Used: false})

	qr, err := store.GetByID("test-id")
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if qr.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if qr.CreatedAt.Before(before) || qr.CreatedAt.After(time.Now()) {
		t.Error("CreatedAt should be approximately now")
	}

	if qr.UsedAt != nil {
		t.Error("UsedAt should be nil before marking as used")
	}

	err = store.MarkAsUsed("test-id")
	if err != nil {
		t.Fatalf("MarkAsUsed() error = %v", err)
	}

	qr, _ = store.GetByID("test-id")
	if qr.UsedAt == nil {
		t.Error("UsedAt should be set after marking as used")
	}
	if !qr.Used {
		t.Error("Used should be true after marking as used")
	}
}
