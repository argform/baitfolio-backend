BEGIN;

ALTER TABLE catches
    DROP CONSTRAINT IF EXISTS tackle_bait_fkey,
    DROP COLUMN IF EXISTS tackle_bait_id,
    ADD COLUMN bait_id SMALLINT,
    ADD COLUMN tackle_id SMALLINT;

DROP TABLE IF EXISTS tackle_bait;

COMMIT;