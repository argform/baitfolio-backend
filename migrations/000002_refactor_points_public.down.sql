BEGIN;

ALTER TABLE points
    DROP CONSTRAINT IF EXISTS points_created_by_fkey;

DROP INDEX IF EXISTS idx_points_created_by;

ALTER TABLE points
    ALTER COLUMN created_by SET NOT NULL;

ALTER TABLE points
    RENAME COLUMN created_by TO owner_id;

ALTER TABLE points
    ADD CONSTRAINT points_owner_id_fkey
    FOREIGN KEY (owner_id)
    REFERENCES users(user_id)
    ON DELETE CASCADE;

CREATE INDEX idx_points_owner_id ON points(owner_id);

ALTER TABLE points
    ADD COLUMN visibility TEXT NOT NULL DEFAULT 'public';

ALTER TABLE points
    ADD CONSTRAINT points_visibility_check
    CHECK (visibility IN ('private', 'public'));

CREATE INDEX idx_points_visibility ON points(visibility);

COMMIT;