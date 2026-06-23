ALTER TABLE users
    ADD COLUMN IF NOT EXISTS token         TEXT,
    ADD COLUMN IF NOT EXISTS refresh_token TEXT;

CREATE INDEX IF NOT EXISTS idx_users_refresh_token ON users (refresh_token)
    WHERE refresh_token IS NOT NULL;
