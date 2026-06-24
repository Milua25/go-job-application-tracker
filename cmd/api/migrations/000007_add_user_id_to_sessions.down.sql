ALTER TABLE sessions
    DROP CONSTRAINT IF EXISTS fk_sessions_user_id;

DROP INDEX IF EXISTS idx_sessions_user_id;

ALTER TABLE sessions
    DROP COLUMN IF EXISTS user_id;