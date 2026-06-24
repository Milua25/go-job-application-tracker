DROP INDEX IF EXISTS idx_sessions_user_email;

ALTER TABLE sessions
    DROP COLUMN IF EXISTS user_email;
