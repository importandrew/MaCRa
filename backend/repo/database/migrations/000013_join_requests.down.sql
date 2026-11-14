DROP TABLE IF EXISTS event_join_requests;
ALTER TABLE events DROP COLUMN IF EXISTS player_limit;
ALTER TABLE events DROP COLUMN IF EXISTS accept_requests;
