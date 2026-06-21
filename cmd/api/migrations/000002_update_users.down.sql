DROP INDEX IF EXISTS idx_users_email_active;
DROP INDEX IF EXISTS idx_users_refresh_token;

ALTER TABLE users
    DROP COLUMN IF EXISTS refresh_token,
    DROP COLUMN IF EXISTS token;
