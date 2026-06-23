DROP INDEX IF EXISTS idx_users_refresh_token;

ALTER TABLE users
    DROP COLUMN IF EXISTS token,
    DROP COLUMN IF EXISTS refresh_token;
