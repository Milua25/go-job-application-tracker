CREATE TABLE IF NOT EXISTS users (
    id                UUID        PRIMARY KEY,
    email             TEXT        NOT NULL UNIQUE,
    password_hash     TEXT        NOT NULL,
    first_name        TEXT        NOT NULL,
    last_name         TEXT        NOT NULL,
    timezone          TEXT        NOT NULL DEFAULT 'UTC',
    is_active         BOOLEAN     NOT NULL DEFAULT TRUE,
    email_verified_at TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);
