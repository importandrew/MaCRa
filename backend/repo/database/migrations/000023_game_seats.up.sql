-- Who sits at a seat (docs/user-scenarios.md, area 8, "Substitutes"). A
-- seat stays its player's (game_players): rounds, scores and results keep
-- their ID. A row here says that from the game's log entry from_log_size
-- on, sitter_id sits there instead (a substitute, or the player again). A
-- seat without rows is sat by its player.
CREATE TABLE game_seats (
    id SERIAL PRIMARY KEY,
    game_id INTEGER NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    seat_player_id INTEGER NOT NULL,
    sitter_id INTEGER NOT NULL REFERENCES users(id),
    from_log_size INTEGER NOT NULL CHECK (from_log_size >= 0),
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (game_id, seat_player_id) REFERENCES game_players(game_id, player_id) ON DELETE CASCADE
);
CREATE INDEX idx_game_seats_game_id ON game_seats(game_id);
CREATE INDEX idx_game_seats_sitter_id ON game_seats(sitter_id);

-- forfeit: the seat's player gets 0 TP and 0 GP for the game (the referee's
-- choice when seating a substitute); the seat's result goes to nobody.
ALTER TABLE game_players ADD COLUMN forfeit BOOLEAN NOT NULL DEFAULT FALSE;

-- game_sitters: who sits at each seat now.
CREATE VIEW game_sitters AS
SELECT gp.game_id, gp.player_id AS seat_player_id,
    COALESCE((SELECT s.sitter_id FROM game_seats s
        WHERE s.game_id = gp.game_id AND s.seat_player_id = gp.player_id
        ORDER BY s.id DESC LIMIT 1), gp.player_id) AS sitter_id
FROM game_players gp;
