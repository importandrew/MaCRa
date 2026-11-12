-- Every time column is timestamptz without a default: the backend sets each
-- value itself, in UTC.

-- The backend always wrote these as UTC wall-clock time.
ALTER TABLE games ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE games
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN finished_at TYPE TIMESTAMPTZ USING finished_at AT TIME ZONE 'UTC';

ALTER TABLE rounds ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE rounds
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';

ALTER TABLE player_results ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE player_results
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';

-- users.created_at was only ever filled by its default, in the server's time
-- zone, which the plain cast assumes too.
ALTER TABLE users ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE users
    ALTER COLUMN created_at TYPE TIMESTAMPTZ;

ALTER TABLE events ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE event_user_roles ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE auth_tokens ALTER COLUMN created_at DROP DEFAULT;
