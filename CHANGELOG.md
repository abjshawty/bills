# Changelog

All notable changes to this project will be documented in this file.

## [0.0.1-alpha] - 2026-03-04

### Added
- `QRCode` struct with `id`, `image`, `client_number`, and `used` fields (JSON and DB tagged)
- `Store` interface defining `Create`, `List`, `GetByID`, and `GetByClientNumber`
- `MemStore` — thread-safe in-memory store implementation using `sync.RWMutex`
- `PostgresStore` — PostgreSQL implementation backed by `database/sql` and `pgx/v5`
- CRUD HTTP endpoints via `net/http` ServeMux with method+path routing (Go 1.22):
  - `POST   /qrcodes` — create a QR code
  - `GET    /qrcodes` — list all QR codes
  - `GET    /qrcodes/{id}` — look up by ID
  - `GET    /qrcodes/phone/{phone}` — look up by client number
- `DATABASE_URL` environment variable for PostgreSQL connection string

### Changed
- Replaced bare-bones `status` handler and default `http.DefaultServeMux` with a structured `Handler` + dedicated `ServeMux`
