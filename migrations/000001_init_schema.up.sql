BEGIN;

CREATE TABLE users (
    user_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    about VARCHAR(250),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE points (
    point_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    lat DOUBLE PRECISION NOT NULL,
    lon DOUBLE PRECISION NOT NULL,
    visibility TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT points_lat_check CHECK (lat BETWEEN -90 AND 90),
    CONSTRAINT points_lon_check CHECK (lon BETWEEN -180 AND 180),
    CONSTRAINT points_visibility_check CHECK (visibility IN ('private', 'public'))
);

CREATE INDEX idx_points_owner_id ON points(owner_id);
CREATE INDEX idx_points_visibility ON points(visibility);

CREATE TABLE reviews (
    review_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_id BIGINT REFERENCES users(user_id) ON DELETE SET NULL,
    point_id BIGINT NOT NULL REFERENCES points(point_id) ON DELETE CASCADE,
    score SMALLINT NOT NULL,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT reviews_score_check CHECK (score BETWEEN 1 AND 5)
);

CREATE UNIQUE INDEX uq_reviews_author_point_not_null
    ON reviews(author_id, point_id)
    WHERE author_id IS NOT NULL;

CREATE INDEX idx_reviews_point_id ON reviews(point_id);
CREATE INDEX idx_reviews_author_id ON reviews(author_id);

CREATE TABLE trips (
    trip_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    point_id BIGINT NOT NULL REFERENCES points(point_id) ON DELETE RESTRICT,
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT trips_time_check CHECK (ended_at IS NULL OR ended_at >= started_at)
);

CREATE INDEX idx_trips_user_id ON trips(user_id);
CREATE INDEX idx_trips_point_id ON trips(point_id);
CREATE INDEX idx_trips_started_at ON trips(started_at);

CREATE TABLE species (
    species_id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(30) UNIQUE NOT NULL
);

CREATE TABLE tackle (
    tackle_id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(60) UNIQUE NOT NULL
);

CREATE TABLE baits (
    bait_id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(60) UNIQUE NOT NULL
);

CREATE TABLE catches (
    catch_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    trip_id BIGINT NOT NULL REFERENCES trips(trip_id) ON DELETE CASCADE,
    species_id SMALLINT REFERENCES species(species_id) ON DELETE SET NULL,
    tackle_id SMALLINT REFERENCES tackle(tackle_id) ON DELETE SET NULL,
    bait_id SMALLINT REFERENCES baits(bait_id) ON DELETE SET NULL,
    weight_grams INTEGER,
    length_cm NUMERIC(6,2),
    released BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    caught_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT catches_weight_check CHECK (weight_grams IS NULL OR weight_grams >= 0),
    CONSTRAINT catches_length_check CHECK (length_cm IS NULL OR length_cm >= 0)
);

CREATE INDEX idx_catches_trip_id ON catches(trip_id);
CREATE INDEX idx_catches_species_id ON catches(species_id);
CREATE INDEX idx_catches_tackle_id ON catches(tackle_id);
CREATE INDEX idx_catches_bait_id ON catches(bait_id);
CREATE INDEX idx_catches_caught_at ON catches(caught_at);

COMMIT;