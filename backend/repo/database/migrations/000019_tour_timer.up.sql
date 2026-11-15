-- A tour's timer: its length, copied from the tournament's settings when it
-- starts (a later change applies to the next tours), and the time a referee
-- added to a single table.
ALTER TABLE tours ADD COLUMN minutes INTEGER CHECK (minutes > 0);
ALTER TABLE games ADD COLUMN extra_minutes INTEGER NOT NULL DEFAULT 0 CHECK (extra_minutes >= 0);
