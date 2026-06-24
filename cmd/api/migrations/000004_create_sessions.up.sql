CREATE TABLE IF NOT EXISTS sessions (
    id            UUID        PRIMARY KEY,
    refresh_token VARCHAR(512) NOT NULL UNIQUE,
    is_revoked    BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at    TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token ON sessions (refresh_token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at  ON sessions (expires_at);
