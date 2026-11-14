-- Event settings (models.EventSettings) as JSON: who may finish or abort a
-- game. Unlike the game rules they can change at any time. Existing events
-- get the defaults, which is how they have worked so far.
ALTER TABLE events ADD COLUMN settings JSONB;
UPDATE events SET settings = '{"finish_game": "seated", "abort_game": "seated_before_first_round"}';
ALTER TABLE events ALTER COLUMN settings SET NOT NULL;
