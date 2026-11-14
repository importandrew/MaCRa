-- Game rules (models.GameRules) as JSON: an event's rules, frozen after its
-- first game, and the copy each game started with and is replayed by.
-- Existing events and games get the MCR standard rules they were played by.
ALTER TABLE events ADD COLUMN rules JSONB;
UPDATE events SET rules = '{"hands": 16, "min_hand": 8}';
ALTER TABLE events ALTER COLUMN rules SET NOT NULL;

ALTER TABLE games ADD COLUMN rules JSONB;
UPDATE games SET rules = '{"hands": 16, "min_hand": 8}';
ALTER TABLE games ALTER COLUMN rules SET NOT NULL;
