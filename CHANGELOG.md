# Changelog

All notable changes to this project will be documented in this file.

## [0.0.1-alpha] - 2026-03-09

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
- `PATCH  /qrcodes/{id}/use` — mark a ticket as used
- `GET    /scan/{id}` — scan alias that also marks a ticket as used
- `GET    /image/{id}` — returns a QR code PNG whose content is the ticket's scan URL
- `GET    /docs/` — Swagger UI
- `GET    /docs/openapi.yaml` — raw OpenAPI spec
- `PORT` and `BASE_URL` environment variables
- UUID v4 auto-generation for ticket IDs server-side
- `Migrate()` on `PostgresStore` — auto-creates the `qrcodes` table on startup
- Dynamic QR code image generation via `github.com/skip2/go-qrcode`

### Changed
- Replaced bare-bones `status` handler and default `http.DefaultServeMux` with a structured `Handler` + dedicated `ServeMux`
