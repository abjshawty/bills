package main

import "errors"

var ErrNotFound = errors.New("not found")

type QRCode struct {
	ID           string `json:"id"            db:"id"`
	Image        string `json:"image"         db:"image"`
	ClientNumber string `json:"client_number" db:"client_number"`
	Used         bool   `json:"used"          db:"used"`
}

type Store interface {
	Create(qr QRCode) error
	List() ([]QRCode, error)
	GetByID(id string) (QRCode, error)
	GetByClientNumber(phone string) (QRCode, error)
}
