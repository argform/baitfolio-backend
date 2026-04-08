BEGIN;

ALTER TABLE points
    DROP CONSTRAINT IF EXISTS points_visibility_check;

DROP INDEX IF EXISTS idx_points_visibility;

ALTER TABLE points
    DROP COLUMN IF EXISTS visibility;

ALTER TABLE points
    DROP CONSTRAINT IF EXISTS points_owner_id_fkey;

ALTER TABLE points
    RENAME COLUMN owner_id TO created_by;

ALTER TABLE points
    ALTER COLUMN created_by DROP NOT NULL;

ALTER TABLE points
    ADD CONSTRAINT points_created_by_fkey
    FOREIGN KEY (created_by)
    REFERENCES users(user_id)
    ON DELETE SET NULL;

DROP INDEX IF EXISTS idx_points_owner_id;
CREATE INDEX idx_points_created_by ON points(created_by);

COMMIT;