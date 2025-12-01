-- Games table - stores game metadata
-- Note: Player positions are stored in game_players table, not here
CREATE TABLE IF NOT EXISTS games(
   id serial PRIMARY KEY,
   status VARCHAR (30) NOT NULL CHECK (status IN ('created', 'live', 'ended', 'aborted')),
   created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   finished_at timestamp
);

-- Index for faster status queries (e.g., finding live games)
CREATE INDEX IF NOT EXISTS idx_games_status ON games(status);
CREATE INDEX IF NOT EXISTS idx_games_created_at ON games(created_at);

