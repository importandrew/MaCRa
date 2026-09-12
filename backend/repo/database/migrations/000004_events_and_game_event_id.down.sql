ALTER TABLE games DROP CONSTRAINT IF EXISTS fk_games_event_id;
DROP INDEX IF EXISTS idx_games_event_id;
DROP INDEX IF EXISTS idx_games_event_id_status;
ALTER TABLE games DROP COLUMN IF EXISTS event_id;

DROP TABLE IF EXISTS event_user_roles;
DROP TABLE IF EXISTS events;
