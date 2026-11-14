-- Requests to join an event: the owner switches them on (public events
-- only) and accepts or declines each one. An accepted request is deleted (the
-- user is a player then); a declined one stays, so the user can't ask again.
ALTER TABLE events ADD COLUMN accept_requests BOOLEAN NOT NULL DEFAULT false;

-- The most players an event takes by a request or the invite link; NULL
-- means no limit. The owner can still add players by hand.
ALTER TABLE events ADD COLUMN player_limit INTEGER CHECK (player_limit > 0);

CREATE TABLE event_join_requests (
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'declined')),
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (event_id, user_id)
);

CREATE INDEX idx_event_join_requests_user_id ON event_join_requests(user_id);
