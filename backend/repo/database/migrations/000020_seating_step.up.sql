-- The step of an interval seating (places i, i+k, i+2k, i+3k at table i);
-- NULL for the other methods.
ALTER TABLE tours ADD COLUMN seating_step INTEGER CHECK (seating_step > 0);
