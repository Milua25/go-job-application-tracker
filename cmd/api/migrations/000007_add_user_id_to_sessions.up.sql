ALTER TABLE sessions
    ADD COLUMN IF NOT EXISTS user_id UUID;

UPDATE sessions s
SET user_id = u.id
FROM users u
WHERE s.user_email = u.email
  AND s.user_id IS NULL;

ALTER TABLE sessions
    ALTER COLUMN user_id SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_sessions_user_id'
    ) THEN
        ALTER TABLE sessions
            ADD CONSTRAINT fk_sessions_user_id
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);