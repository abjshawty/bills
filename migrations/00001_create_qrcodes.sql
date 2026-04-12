-- +goose Up
CREATE TABLE IF NOT EXISTS qrcodes (
    id            TEXT PRIMARY KEY,
    image         TEXT        NOT NULL,
    client_number TEXT        NOT NULL UNIQUE,
    used          BOOLEAN     NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    used_at       TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS qrcodes;
