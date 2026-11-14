-- A tour's seating is first generated as a draft, seen only by the
-- tournament's staff, then published to the players.
ALTER TABLE tours DROP CONSTRAINT tours_status_check;
ALTER TABLE tours ADD CONSTRAINT tours_status_check
    CHECK (status IN ('draft', 'published', 'live', 'finished', 'aborted'));
ALTER TABLE tours ALTER COLUMN published_at DROP NOT NULL;
