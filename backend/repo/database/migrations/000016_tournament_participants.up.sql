-- A tournament's participant numbers, from the referee's draw, and who
-- withdrew. The numbers are valid while they belong to exactly the
-- tournament's current players; after a change the draw is redone.
CREATE TABLE tournament_participants (
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    number INTEGER NOT NULL CHECK (number > 0),
    -- When the player withdrew (after the first tour started); they keep
    -- their past games and their number, and aren't seated any more.
    withdrawn_at TIMESTAMPTZ,
    PRIMARY KEY (event_id, user_id),
    UNIQUE (event_id, number)
);
