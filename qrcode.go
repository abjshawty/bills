package main

import (
	"errors"
	"time"
)

// ErrNotFound is returned by Store methods when a QR code does not exist.
var ErrNotFound = errors.New("not found")

// ErrAlreadyUsed is returned when attempting to mark an already-used ticket.
var ErrAlreadyUsed = errors.New("already used")

// QRCode represents a single event ticket backed by a QR code image.
type QRCode struct {
	ID           string     `json:"id"            db:"id"`
	Image        string     `json:"image"         db:"image"`
	ClientNumber string     `json:"client_number" db:"client_number"`
	Used         bool       `json:"used"          db:"used"`
	CreatedAt    time.Time  `json:"created_at"    db:"created_at"`
	UsedAt       *time.Time `json:"used_at"       db:"used_at"`
}

// Store is the persistence interface for QR code tickets.
// Any backend (Postgres, in-memory, …) must satisfy this interface.
type Store interface {
	// Create persists a new QR code ticket.
	Create(qr QRCode) error
	// List returns all stored QR code tickets.
	List() ([]QRCode, error)
	// GetByID returns the ticket with the given ID, or ErrNotFound.
	GetByID(id string) (QRCode, error)
	// GetByClientNumber returns all tickets for the given phone number, or ErrNotFound.
	GetByClientNumber(phone string) ([]QRCode, error)
	// MarkAsUsed marks the ticket with the given ID as used.
	MarkAsUsed(id string) error
	// Delete removes the ticket with the given ID.
	Delete(id string) error
}
