-- Who entered each round log entry; NULL for entries saved before this
-- migration (and for a deleted account).
ALTER TABLE rounds ADD COLUMN entered_by INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- Undone round log entries. An undo still deletes the entry from rounds (the
-- game is replayed from rounds alone); a copy of it lands here with who undid
-- it and when, so everyone who sees the game's log sees the undo too.
CREATE TABLE undone_rounds (
    id SERIAL PRIMARY KEY,
    game_id INTEGER NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    round_number INTEGER NOT NULL,
    outcome VARCHAR(20) NOT NULL CHECK (outcome IN ('deal-in', 'self', 'draw', 'false mahjong', 'penalty')),
    winner_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    loser_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    value INTEGER,
    force_next_round BOOLEAN NOT NULL,
    entered_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    entered_at TIMESTAMPTZ NOT NULL,
    undone_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    undone_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_undone_rounds_game_id ON undone_rounds(game_id);
