-- One setting decides who sees an event: its visibility. A public event
-- shows everything but managing it to everyone, a private one to its
-- members; the separate reaches of the rating and the games are gone
-- (docs/user-scenarios.md, area 1).
ALTER TABLE events DROP COLUMN ratings_reach;
ALTER TABLE events DROP COLUMN games_reach;
