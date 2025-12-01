-- Game players junction table
-- Links players to games with their starting wind positions
CREATE TABLE IF NOT EXISTS game_players (
    id serial PRIMARY KEY,
    game_id integer NOT NULL,
    player_id integer NOT NULL,
    starting_wind varchar(10) NOT NULL CHECK (starting_wind IN ('east', 'south', 'west', 'north')),
    UNIQUE(game_id, player_id),
    UNIQUE(game_id, starting_wind),
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Rounds table - stores all round events (outcomes and penalties) across all games
-- Multiple events can occur in the same round (e.g., penalties that don't advance the round)
-- MCR games always consist of exactly 16 rounds
CREATE TABLE IF NOT EXISTS rounds (
    id serial PRIMARY KEY,
    game_id integer NOT NULL,
    round_number integer NOT NULL,
    event_sequence integer NOT NULL,
    outcome varchar(20) NOT NULL CHECK (outcome IN ('deal-in', 'self', 'draw', 'false mahjong', 'penalty')),
    winner_id integer,
    loser_id integer,
    value integer,
    force_next_round boolean NOT NULL DEFAULT false,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    FOREIGN KEY (winner_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (loser_id) REFERENCES users(id) ON DELETE SET NULL,
    CHECK (round_number >= 1 AND round_number <= 16),
    CHECK (event_sequence > 0),
    UNIQUE(game_id, round_number, event_sequence)
);

-- Function to auto-increment event_sequence for each (game_id, round_number) combination
CREATE OR REPLACE FUNCTION set_event_sequence()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.event_sequence IS NULL THEN
        SELECT COALESCE(MAX(event_sequence), 0) + 1
        INTO NEW.event_sequence
        FROM rounds
        WHERE game_id = NEW.game_id AND round_number = NEW.round_number;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically set event_sequence before insert
CREATE TRIGGER set_event_sequence_trigger
    BEFORE INSERT ON rounds
    FOR EACH ROW
    EXECUTE FUNCTION set_event_sequence();

-- Round scores table - stores player scores AFTER each round event
-- Each entry in rounds table has corresponding score entries showing the state after that event
CREATE TABLE IF NOT EXISTS round_scores (
    id serial PRIMARY KEY,
    round_id integer NOT NULL,
    game_id integer NOT NULL,
    round_number integer NOT NULL,
    event_sequence integer NOT NULL,
    player_id integer NOT NULL,
    score integer NOT NULL,
    FOREIGN KEY (round_id) REFERENCES rounds(id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(round_id, player_id),
    CHECK (round_number >= 1 AND round_number <= 16),
    CHECK (event_sequence > 0)
);

-- Player results table - stores individual player results for completed games
-- Each game should have exactly 4 player results (one per player)
CREATE TABLE IF NOT EXISTS player_results (
    id serial PRIMARY KEY,
    game_id integer NOT NULL,
    player_id integer NOT NULL,
    game_points integer NOT NULL,
    tournament_points numeric(4,3) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(game_id, player_id)
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_game_players_game_id ON game_players(game_id);
CREATE INDEX IF NOT EXISTS idx_game_players_player_id ON game_players(player_id);
CREATE INDEX IF NOT EXISTS idx_rounds_game_id ON rounds(game_id);
CREATE INDEX IF NOT EXISTS idx_rounds_round_number ON rounds(game_id, round_number);
CREATE INDEX IF NOT EXISTS idx_rounds_event_sequence ON rounds(game_id, round_number, event_sequence);
CREATE INDEX IF NOT EXISTS idx_round_scores_round_id ON round_scores(round_id);
CREATE INDEX IF NOT EXISTS idx_round_scores_game_id ON round_scores(game_id);
CREATE INDEX IF NOT EXISTS idx_round_scores_event ON round_scores(game_id, round_number, event_sequence);
CREATE INDEX IF NOT EXISTS idx_round_scores_player_id ON round_scores(player_id);
CREATE INDEX IF NOT EXISTS idx_player_results_game_id ON player_results(game_id);
CREATE INDEX IF NOT EXISTS idx_player_results_player_id ON player_results(player_id);

-- Prevent duplicate main outcomes (deal-in, self, draw) per round
-- These outcomes are mutually exclusive and can only happen once per round
-- Penalties can still occur multiple times in the same round
CREATE UNIQUE INDEX IF NOT EXISTS idx_rounds_unique_main_outcome
ON rounds(game_id, round_number, outcome)
WHERE outcome IN ('deal-in', 'self', 'draw');

