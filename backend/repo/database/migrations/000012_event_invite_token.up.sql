-- An event's invite link: whoever opens it while signed in becomes a player.
-- NULL means the link is off. The token is stored as it is: the owner sees
-- the link again in Manage.
ALTER TABLE events ADD COLUMN invite_token TEXT;
CREATE UNIQUE INDEX events_invite_token_key ON events (invite_token) WHERE invite_token IS NOT NULL;
