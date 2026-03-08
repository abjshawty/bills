# CLAUDE.md

## Project Overview
Go HTTP server for managing QR code tickets for a livestream event. Module name: `tickets-go`.

## Architecture

| File | Role |
|------|------|
| `qrcode.go` | `QRCode` struct, `Store` interface, `ErrNotFound` sentinel |
| `store.go` | `MemStore` — in-memory `Store` implementation (useful for tests) |
| `db.go` | `PostgresStore` — Postgres `Store` implementation via `pgx/stdlib` |
| `handlers.go` | HTTP handlers wired to a `Store`; auto-generates UUIDs on create |
| `docs.go` | Embedded OpenAPI spec + Swagger UI served at `/docs/` |
| `main.go` | Entry point: loads `.env`, connects to Postgres, starts server |

## Running Locally
```bash
cp .env.example .env   # fill in DATABASE_URL
go run .
```
Server listens on `:9000`.

## Environment Variables
| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | yes | — | Postgres connection string |
| `PORT` | no | `9000` | Port the server listens on |

`.env` is loaded automatically via `godotenv`. Real env vars take precedence.

## API Endpoints
| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/qrcodes` | Create a ticket (ID auto-generated) |
| `GET` | `/qrcodes` | List all tickets |
| `GET` | `/qrcodes/{id}` | Get by ID |
| `GET` | `/qrcodes/phone/{phone}` | Get by client phone number |
| `PATCH` | `/qrcodes/{id}/use` | Mark a ticket as used |
| `GET` | `/docs/` | Swagger UI |
| `GET` | `/docs/openapi.yaml` | Raw OpenAPI spec |

## Key Conventions
- Store is accessed only through the `Store` interface — keep handlers decoupled from the backend.
- Use `ErrNotFound` (not raw `sql.ErrNoRows`) when returning not-found errors from any store implementation.
- IDs are UUID v4 strings generated server-side in the `Create` handler — never trust client-supplied IDs.
- No test suite yet — prefer table-driven tests with `MemStore` as the backend if added.

## Dependencies
- `github.com/jackc/pgx/v5` — Postgres driver (used via `database/sql` stdlib adapter)
- `github.com/joho/godotenv` — `.env` file loading
- `github.com/google/uuid` — UUID v4 generation for ticket IDs
