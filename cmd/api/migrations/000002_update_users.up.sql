ALTER TABLE users
    ADD COLUMN IF NOT EXISTS token         TEXT,
    ADD COLUMN IF NOT EXISTS refresh_token TEXT;

-- Fast lookup when validating a refresh token.
CREATE INDEX IF NOT EXISTS idx_users_refresh_token ON users (refresh_token)
    WHERE refresh_token IS NOT NULL;

-- Enforce uniqueness only on active (non-deleted) rows so soft-deleted
-- emails can be re-registered.
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active ON users (email)
    WHERE deleted_at IS NULL;
