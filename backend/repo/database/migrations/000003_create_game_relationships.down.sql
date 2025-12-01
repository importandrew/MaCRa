-- Rollback game relationships tables
DROP INDEX IF EXISTS idx_rounds_unique_main_outcome;
DROP INDEX IF EXISTS idx_player_results_player_id;
DROP INDEX IF EXISTS idx_player_results_game_id;
DROP INDEX IF EXISTS idx_round_scores_player_id;
DROP INDEX IF EXISTS idx_round_scores_event;
DROP INDEX IF EXISTS idx_round_scores_game_id;
DROP INDEX IF EXISTS idx_round_scores_round_id;
DROP INDEX IF EXISTS idx_rounds_round_number;
DROP INDEX IF EXISTS idx_rounds_game_id;
DROP INDEX IF EXISTS idx_game_players_player_id;
DROP INDEX IF EXISTS idx_game_players_game_id;

DROP TRIGGER IF EXISTS set_event_sequence_trigger ON rounds;
DROP FUNCTION IF EXISTS set_event_sequence();

DROP TABLE IF EXISTS player_results;
DROP TABLE IF EXISTS round_scores;
DROP TABLE IF EXISTS rounds;
DROP TABLE IF EXISTS game_players;

