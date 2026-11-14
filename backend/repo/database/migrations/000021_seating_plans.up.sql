-- A tournament's seating prepared for all tours at once, generated or
-- imported from a file: per tour (the JSON object's keys), its tables, each
-- four participant numbers east to north. A tour seated "fixed" takes its
-- tables from here.
CREATE TABLE seating_plans (
    event_id INTEGER PRIMARY KEY REFERENCES events(id) ON DELETE CASCADE,
    source VARCHAR(20) NOT NULL CHECK (source IN ('generated', 'imported')),
    created_at TIMESTAMPTZ NOT NULL,
    tours JSONB NOT NULL
);
