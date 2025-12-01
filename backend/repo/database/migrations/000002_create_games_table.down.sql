-- Rollback games table
DROP INDEX IF EXISTS idx_games_created_at;
DROP INDEX IF EXISTS idx_games_status;
DROP TABLE IF EXISTS games;

