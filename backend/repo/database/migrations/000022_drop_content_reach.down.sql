ALTER TABLE events ADD COLUMN ratings_reach VARCHAR(20) NOT NULL DEFAULT 'members'
    CHECK (ratings_reach IN ('public', 'members', 'staff'));
ALTER TABLE events ADD COLUMN games_reach VARCHAR(20) NOT NULL DEFAULT 'members'
    CHECK (games_reach IN ('public', 'members', 'staff'));
UPDATE events SET ratings_reach = 'public', games_reach = 'public' WHERE visibility = 'public';
