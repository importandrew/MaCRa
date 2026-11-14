DROP INDEX IF EXISTS events_invite_token_key;
ALTER TABLE events DROP COLUMN IF EXISTS invite_token;
