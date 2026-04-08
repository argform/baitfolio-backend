BEGIN;

ALTER TABLE points
    DROP CONSTRAINT IF EXISTS points_waterbody_hydrology_fkey,
    DROP CONSTRAINT IF EXISTS points_shore_type_fkey,
    DROP CONSTRAINT IF EXISTS points_access_type_fkey,
    DROP COLUMN IF EXISTS waterbody_hydrology_id,
    DROP COLUMN IF EXISTS shore_type_id,
    DROP COLUMN IF EXISTS access_type_id;

DROP TABLE IF EXISTS waterbody_hydrology;
DROP TABLE IF EXISTS access_type;
DROP TABLE IF EXISTS hydrology_type;
DROP TABLE IF EXISTS waterbody_type;
DROP TABLE IF EXISTS shore_type;

COMMIT;