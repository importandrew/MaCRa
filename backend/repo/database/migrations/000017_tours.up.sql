-- A tournament's tours. Publishing a tour's seating creates its games
-- (status created); cancelling the seating deletes the tour with its games.
-- An aborted tour stays with its aborted games and the number is played
-- again, so the number is unique only among the tours not aborted.
CREATE TABLE tours (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    number INTEGER NOT NULL CHECK (number > 0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('published', 'live', 'finished', 'aborted')),
    seating VARCHAR(20) NOT NULL,
    published_at TIMESTAMPTZ NOT NULL,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX tours_event_number_key ON tours(event_id, number) WHERE status <> 'aborted';
CREATE INDEX idx_tours_event_id ON tours(event_id);

-- A tournament game belongs to a tour and a table.
ALTER TABLE games ADD COLUMN tour_id INTEGER REFERENCES tours(id) ON DELETE CASCADE;
ALTER TABLE games ADD COLUMN table_number INTEGER CHECK (table_number > 0);
CREATE INDEX idx_games_tour_id ON games(tour_id);
