# TODO

## UX / Responses
- [ ] Better view results on scans (e.g. styled HTML response instead of raw JSON)
- [ ] Image downloads (e.g. `Content-Disposition: attachment` header on `GET /image/{id}`)

## Validation & Error Handling
- [ ] Validate `client_number` format on `POST /qrcodes` (e.g. reject empty or non-numeric strings)
- [ ] Return `409 Conflict` on `POST /qrcodes` when `client_number` already exists (currently a raw 500 from the DB unique constraint)
- [ ] Guard `PATCH /qrcodes/{id}/use` against double-use (return `409` if ticket is already used)

## Data Model
- [ ] Add a `created_at` timestamp field to `QRCode` and the `qrcodes` table
- [ ] Add a `used_at` nullable timestamp that is set when `MarkAsUsed` is called

## API / Routes
- [ ] `DELETE /qrcodes/{id}` — hard-delete a ticket
- [ ] `GET /scan/{id}` currently shares `MarkAsUsed` handler — give it its own handler that renders the styled HTML scan page (see UX item above)
- [ ] Update OpenAPI spec (`docs.go`) to document `GET /image/{id}`, `GET /scan/{id}`, `created_at`, and `used_at` fields

## README
- [ ] Fix `.env.example` reference in README (file is named `.env.sample` in the repo)
- [ ] Document `BASE_URL` env variable in the environment variables table
- [ ] Add `GET /image/{id}` and `GET /scan/{id}` to the endpoints table

## Testing
- [ ] Unit tests for `MemStore` (all five methods)
- [ ] HTTP handler tests using `httptest` (create, list, get, mark-as-used, image, scan)

## Ops / Quality
- [ ] Add a `Dockerfile` for containerised deployment
- [ ] Structured logging (e.g. `log/slog`) instead of bare `log.Fatal` calls in `main.go`
- [ ] Graceful shutdown (`signal.NotifyContext` + `http.Server.Shutdown`)
