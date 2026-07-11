CREATE EXTENSION IF NOT EXISTS citext;

DO $$ BEGIN
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'media_type') THEN
        CREATE TYPE media_type AS ENUM (
            'FILM',
            'BOOK',
            'GAME',
            'SHOW',
            'SEASON',
            'EPISODE',
            'ALBUM',
            'SONG'
        );
    END IF;

    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'film_type') THEN
        CREATE TYPE film_type AS ENUM ('FEATURE FILM', 'SHORT FILM');
    END IF;
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'book_type') THEN
        CREATE TYPE book_type AS ENUM ('NOVEL', 'SHORT STORY', 'COLLECTION');
    END IF;
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'game_type') THEN
        CREATE TYPE game_type AS ENUM ('VIDEO GAME');
    END IF;
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'show_type') THEN
        CREATE TYPE show_type AS ENUM ('SERIES', 'MINI SERIES');
    END IF;
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'music_type') THEN
        CREATE TYPE music_type AS ENUM ('ALBUM', 'EXTENDED PLAY', 'MIXTAPE');
    END IF;
END $$;



CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    email CITEXT UNIQUE NOT NULL,
    display_name CITEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    last_login TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS films(
    id BIGSERIAL PRIMARY KEY,
    film_type film_type NOT NULL,
    title CITEXT NOT NULL,
    description TEXT,
    release_date DATE NOT NULL,
    runtime_minutes INTEGER,
    external_references JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS books(
    id BIGSERIAL PRIMARY KEY,
    book_type book_type NOT NULL,
    title CITEXT NOT NULL,
    description TEXT,
    release_date DATE NOT NULL,
    external_references JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS games(
    id BIGSERIAL PRIMARY KEY,
    game_type game_type NOT NULL,
    title CITEXT NOT NULL,
    description TEXT,
    release_date DATE NOT NULL,
    external_references JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS dlc(
    id BIGSERIAL PRIMARY KEY,
    game_id BIGINT REFERENCES games(id) ON DELETE CASCADE,
    title CITEXT NOT NULL,
    description TEXT,
    release_date DATE NOT NULL,
    external_references JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS shows(
    id BIGSERIAL PRIMARY KEY,
    show_type show_type NOT NULL,
    title CITEXT NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE,
    external_references JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS seasons(
    id BIGSERIAL PRIMARY KEY,
    show_id BIGINT REFERENCES shows(id) ON DELETE CASCADE,
    season_number INTEGER NOT NULL,
    title TEXT,
    start_date DATE NOT NULL,
    end_date DATE,
    external_references JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS episodes(
    id BIGSERIAL PRIMARY KEY,
    season_id BIGINT REFERENCES seasons(id) ON DELETE CASCADE,
    episode_number INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    release_date DATE NOT NULL,
    runtime_minutes INTEGER,
    external_references JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS albums(
    id BIGSERIAL PRIMARY KEY,
    music_type music_type NOT NULL,
    title CITEXT NOT NULL,
    description TEXT,
    release_date DATE NOT NULL,
    external_references JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS songs(
    id BIGSERIAL PRIMARY KEY,
    album_id BIGINT REFERENCES albums(id) ON DELETE CASCADE,
    song_number INTEGER NOT NULL,
    title TEXT,
    runtime_seconds INTEGER,
    external_references JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS user_media(
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    media_id BIGINT,
    media_type media_type NOT NULL,
    rating INTEGER CHECK (rating BETWEEN 1 AND 100),
    PRIMARY KEY (user_id, media_id, media_type)
);

CREATE TABLE IF NOT EXISTS diary_entry(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    media_id BIGINT,
    media_type media_type NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    rating INTEGER CHECK (rating BETWEEN 1 AND 100),
    description TEXT,
    is_private BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON films
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON books
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON games
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON dlc
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON shows
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON seasons
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON episodes
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON albums
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON songs
FOR EACH ROW EXECUTE FUNCTION set_updated_at();