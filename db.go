package main

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// ErrDuplicateClientNumber is returned when attempting to create a ticket with an existing client_number.
var ErrDuplicateClientNumber = errors.New("client number already exists")

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

// Migrate runs all pending SQL migrations from the embedded migrations directory.
func (s *PostgresStore) Migrate() error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose set dialect: %w", err)
	}
	if err := goose.Up(s.db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

func (s *PostgresStore) Create(qr QRCode) error {
	_, err := s.db.Exec(
		`INSERT INTO qrcodes (id, image, client_number, used, created_at) VALUES ($1, $2, $3, $4, $5)`,
		qr.ID, qr.Image, qr.ClientNumber, qr.Used, qr.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return ErrDuplicateClientNumber
		}
		return err
	}
	return nil
}

func (s *PostgresStore) List() ([]QRCode, error) {
	rows, err := s.db.Query(`SELECT id, image, client_number, used, created_at, used_at FROM qrcodes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []QRCode
	for rows.Next() {
		var qr QRCode
		if err := rows.Scan(&qr.ID, &qr.Image, &qr.ClientNumber, &qr.Used, &qr.CreatedAt, &qr.UsedAt); err != nil {
			return nil, err
		}
		list = append(list, qr)
	}
	return list, rows.Err()
}

func (s *PostgresStore) GetByID(id string) (QRCode, error) {
	var qr QRCode
	err := s.db.QueryRow(
		`SELECT id, image, client_number, used, created_at, used_at FROM qrcodes WHERE id = $1`, id,
	).Scan(&qr.ID, &qr.Image, &qr.ClientNumber, &qr.Used, &qr.CreatedAt, &qr.UsedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return QRCode{}, ErrNotFound
	}
	return qr, err
}

func (s *PostgresStore) GetByClientNumber(phone string) ([]QRCode, error) {
	rows, err := s.db.Query(
		`SELECT id, image, client_number, used, created_at, used_at FROM qrcodes WHERE client_number = $1`, phone,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []QRCode
	for rows.Next() {
		var qr QRCode
		if err := rows.Scan(&qr.ID, &qr.Image, &qr.ClientNumber, &qr.Used, &qr.CreatedAt, &qr.UsedAt); err != nil {
			return nil, err
		}
		list = append(list, qr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, ErrNotFound
	}
	return list, nil
}

// MarkAsUsed sets the used flag to true for the ticket with the given ID.
func (s *PostgresStore) MarkAsUsed(id string) error {
	qr, err := s.GetByID(id)
	if err != nil {
		return err
	}
	if qr.Used {
		return ErrAlreadyUsed
	}

	res, err := s.db.Exec(`UPDATE qrcodes SET used = true, used_at = NOW() WHERE id = $1`, id)
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

// Delete removes the ticket with the given ID.
func (s *PostgresStore) Delete(id string) error {
	res, err := s.db.Exec(`DELETE FROM qrcodes WHERE id = $1`, id)
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
