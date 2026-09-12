-- Events (competition/session).
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL CHECK (type IN ('tournament', 'regular')),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_user_id INTEGER NOT NULL REFERENCES users(id),
    visibility VARCHAR(20) NOT NULL DEFAULT 'private' CHECK (visibility IN ('private', 'public')),
    ratings_reach VARCHAR(20) NOT NULL DEFAULT 'members' CHECK (ratings_reach IN ('public', 'members', 'staff')),
    games_reach VARCHAR(20) NOT NULL DEFAULT 'members' CHECK (games_reach IN ('public', 'members', 'staff')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_events_owner_user_id ON events(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_events_visibility ON events(visibility);

CREATE TABLE IF NOT EXISTS event_user_roles (
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('player', 'referee', 'observer')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id, role)
);
-- Regular-event owner player+referee rows are inserted at create time (not stored as 'owner' role).

CREATE INDEX IF NOT EXISTS idx_event_user_roles_event_id ON event_user_roles(event_id);
CREATE INDEX IF NOT EXISTS idx_event_user_roles_user_id ON event_user_roles(user_id);

ALTER TABLE games ADD COLUMN IF NOT EXISTS event_id INTEGER NOT NULL;

ALTER TABLE games
    ADD CONSTRAINT fk_games_event_id FOREIGN KEY (event_id) REFERENCES events(id);

CREATE INDEX IF NOT EXISTS idx_games_event_id ON games(event_id);
CREATE INDEX IF NOT EXISTS idx_games_event_id_status ON games(event_id, status);