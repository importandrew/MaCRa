DELETE FROM tours WHERE status = 'draft';
ALTER TABLE tours ALTER COLUMN published_at SET NOT NULL;
ALTER TABLE tours DROP CONSTRAINT tours_status_check;
ALTER TABLE tours ADD CONSTRAINT tours_status_check
    CHECK (status IN ('published', 'live', 'finished', 'aborted'));
