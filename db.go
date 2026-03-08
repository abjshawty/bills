package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostgresStore is a Store backed by a PostgreSQL database.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore opens and pings a Postgres connection using connStr.
func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Create(qr QRCode) error {
	_, err := s.db.Exec(
		`INSERT INTO qrcodes (id, image, client_number, used) VALUES ($1, $2, $3, $4)`,
		qr.ID, qr.Image, qr.ClientNumber, qr.Used,
	)
	return err
}

func (s *PostgresStore) List() ([]QRCode, error) {
	rows, err := s.db.Query(`SELECT id, image, client_number, used FROM qrcodes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []QRCode
	for rows.Next() {
		var qr QRCode
		if err := rows.Scan(&qr.ID, &qr.Image, &qr.ClientNumber, &qr.Used); err != nil {
			return nil, err
		}
		list = append(list, qr)
	}
	return list, rows.Err()
}

func (s *PostgresStore) GetByID(id string) (QRCode, error) {
	var qr QRCode
	err := s.db.QueryRow(
		`SELECT id, image, client_number, used FROM qrcodes WHERE id = $1`, id,
	).Scan(&qr.ID, &qr.Image, &qr.ClientNumber, &qr.Used)
	if errors.Is(err, sql.ErrNoRows) {
		return QRCode{}, ErrNotFound
	}
	return qr, err
}

func (s *PostgresStore) GetByClientNumber(phone string) (QRCode, error) {
	var qr QRCode
	err := s.db.QueryRow(
		`SELECT id, image, client_number, used FROM qrcodes WHERE client_number = $1`, phone,
	).Scan(&qr.ID, &qr.Image, &qr.ClientNumber, &qr.Used)
	if errors.Is(err, sql.ErrNoRows) {
		return QRCode{}, ErrNotFound
	}
	return qr, err
}

// MarkAsUsed sets the used flag to true for the ticket with the given ID.
func (s *PostgresStore) MarkAsUsed(id string) error {
	res, err := s.db.Exec(`UPDATE qrcodes SET used = true WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
