-- A Swiss or interval seating may avoid repeat meetings or not; the Swiss
-- seating always did, the interval one never did.
ALTER TABLE tours ADD COLUMN avoid_repeats BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE tours SET avoid_repeats = TRUE WHERE seating = 'swiss';

-- A tournament's plan: per tour, the seating method it is to have, and the
-- tables imported from a file for it (participant numbers, east to north),
-- which the "fixed" method takes. A tour without a row is random.
CREATE TABLE tour_plans (
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    number INTEGER NOT NULL CHECK (number > 0),
    seating VARCHAR(20) NOT NULL,
    seating_step INTEGER CHECK (seating_step > 0),
    avoid_repeats BOOLEAN NOT NULL,
    tables JSONB,
    PRIMARY KEY (event_id, number)
);
