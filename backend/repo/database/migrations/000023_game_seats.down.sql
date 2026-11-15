DROP VIEW IF EXISTS game_sitters;
ALTER TABLE game_players DROP COLUMN IF EXISTS forfeit;
DROP TABLE IF EXISTS game_seats;
