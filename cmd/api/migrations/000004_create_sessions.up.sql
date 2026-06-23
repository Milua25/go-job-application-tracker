CREATE TABLE IF NOT EXISTS sessions (
    id            UUID        PRIMARY KEY,
    user_email    TEXT        NOT NULL REFERENCES users(email) ON DELETE CASCADE,
    token         TEXT        NOT NULL UNIQUE,
    refresh_token TEXT        NOT NULL UNIQUE,
    is_revoked    BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at    TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_email  ON sessions (user_email);
CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token ON sessions (refresh_token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at  ON sessions (expires_at);
