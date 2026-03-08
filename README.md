# Elementary Server for ticket handling

Hi nerds! Last week, i had the genius idea to organize a livestream event for no fucking reason.

Now, i need a system to handle the tickets, and yeah if you go through my github you'll know i already have at least two versions of exactly that. Buttttt those are in Typescript, and i'm learning Go rn. So, i'm following the HTTP Server section of Go By Example, and i'll adapt it to work for my use case. Feel free to comment!

xoxo <3

Timmy

---

## Getting started

### Prerequisites
- Go 1.24+
- A PostgreSQL database (the project uses [Neon](https://neon.tech) by default)

### Setup

```bash
git clone <repo>
cd bills
cp .env.example .env   # fill in your DATABASE_URL
go run .
```

The server starts on the port defined in `.env` (default `9000`).

### Environment variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | yes | — | Postgres connection string |
| `PORT` | no | `9000` | Port the server listens on |

Create a `.env` file in the project root (it is git-ignored):

```env
DATABASE_URL=postgres://user:password@localhost:5432/dbname
PORT=9000
```

---

## API

Interactive docs are available at **`/docs/`** when the server is running.
The raw OpenAPI spec is at **`/docs/openapi.yaml`**.

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/qrcodes` | Create a ticket (ID is auto-generated) |
| `GET` | `/qrcodes` | List all tickets |
| `GET` | `/qrcodes/{id}` | Get a ticket by ID |
| `GET` | `/qrcodes/phone/{phone}` | Get a ticket by client phone number |
| `PATCH` | `/qrcodes/{id}/use` | Mark a ticket as used |

### Example: create a ticket

```bash
curl -X POST http://localhost:9000/qrcodes \
  -H "Content-Type: application/json" \
  -d '{"image": "<base64>", "client_number": "+213XXXXXXXXX"}'
```

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "image": "<base64>",
  "client_number": "+213XXXXXXXXX",
  "used": false
}
```

### Example: mark a ticket as used

```bash
curl -X PATCH http://localhost:9000/qrcodes/550e8400-e29b-41d4-a716-446655440000/use
```

---

## Database schema

```sql
CREATE TABLE qrcodes (
  id            TEXT PRIMARY KEY,
  image         TEXT NOT NULL,
  client_number TEXT NOT NULL UNIQUE,
  used          BOOLEAN NOT NULL DEFAULT false
);
```
