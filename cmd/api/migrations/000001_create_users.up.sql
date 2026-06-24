CREATE TABLE IF NOT EXISTS users (
    id                UUID        PRIMARY KEY,
    email             VARCHAR(254) NOT NULL UNIQUE,
    password_hash     TEXT         NOT NULL,
    first_name        VARCHAR(50)  NOT NULL,
    last_name         VARCHAR(50)  NOT NULL,
    timezone          TEXT        NOT NULL DEFAULT 'UTC',
    is_active         BOOLEAN     NOT NULL DEFAULT TRUE,
    email_verified_at TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);
