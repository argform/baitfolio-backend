BEGIN;

CREATE TABLE tackle_bait (
    tackle_bait_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tackle_id SMALLINT NOT NULL REFERENCES tackle(tackle_id) ON DELETE CASCADE,
    bait_id SMALLINT NOT NULL REFERENCES baits(bait_id) ON DELETE CASCADE,
    CONSTRAINT tackle_bait_unique UNIQUE (tackle_id, bait_id)
);

ALTER TABLE catches
    DROP COLUMN IF EXISTS bait_id,
    DROP COLUMN IF EXISTS tackle_id,
    ADD COLUMN tackle_bait_id INTEGER,
    ADD CONSTRAINT tackle_bait_fkey
        FOREIGN KEY (tackle_bait_id)
        REFERENCES tackle_bait(tackle_bait_id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE;

COMMIT;