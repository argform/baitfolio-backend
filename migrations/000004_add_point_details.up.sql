BEGIN;

CREATE TABLE shore_type (
    shore_type_id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(60) NOT NULL UNIQUE
);

CREATE TABLE waterbody_type (
    waterbody_type_id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(60) NOT NULL UNIQUE
);

CREATE TABLE hydrology_type (
    hydrology_type_id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(60) NOT NULL UNIQUE
);

CREATE TABLE access_type (
    access_type_id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(60) NOT NULL UNIQUE
);

CREATE TABLE waterbody_hydrology (
    waterbody_hydrology_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    waterbody_type_id SMALLINT NOT NULL
        REFERENCES waterbody_type(waterbody_type_id)
        ON DELETE CASCADE,
    hydrology_type_id SMALLINT NOT NULL
        REFERENCES hydrology_type(hydrology_type_id)
        ON DELETE CASCADE,
    CONSTRAINT waterbody_hydrology_unique UNIQUE (waterbody_type_id, hydrology_type_id)
);

ALTER TABLE points
    ADD COLUMN waterbody_hydrology_id INTEGER,
    ADD CONSTRAINT points_waterbody_hydrology_fkey
        FOREIGN KEY (waterbody_hydrology_id)
        REFERENCES waterbody_hydrology(waterbody_hydrology_id)
        ON DELETE SET NULL,
    ADD COLUMN shore_type_id SMALLINT,
    ADD CONSTRAINT points_shore_type_fkey
        FOREIGN KEY (shore_type_id)
        REFERENCES shore_type(shore_type_id)
        ON DELETE SET NULL,
    ADD COLUMN access_type_id SMALLINT,
    ADD CONSTRAINT points_access_type_fkey
        FOREIGN KEY (access_type_id)
        REFERENCES access_type(access_type_id)
        ON DELETE SET NULL;

COMMIT;