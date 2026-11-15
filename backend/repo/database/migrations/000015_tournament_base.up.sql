-- Tournaments, first part: the substitute role, when the first tour started,
-- and the tournament's own settings.

-- A substitute sits only where a referee seats them; a role of its own,
-- never together with player (checked by the backend).
ALTER TABLE event_user_roles DROP CONSTRAINT event_user_roles_role_check;
ALTER TABLE event_user_roles ADD CONSTRAINT event_user_roles_role_check
    CHECK (role IN ('player', 'referee', 'observer', 'substitute'));

-- When the tournament's first tour started; NULL while it is being prepared.
-- With finished_at it gives the status: preparation, in progress, finished.
ALTER TABLE events ADD COLUMN started_at TIMESTAMPTZ;

-- A tournament's settings (number of tours, tour length, what happens when
-- time is up), complete JSON as with rules and settings; NULL for a regular
-- event.
ALTER TABLE events ADD COLUMN tournament JSONB;
UPDATE events
SET tournament = '{"tours": 5, "tour_minutes": 120, "time_up": "finish_hand", "extra_hands": 0}'
WHERE type = 'tournament';
