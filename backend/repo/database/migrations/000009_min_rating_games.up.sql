-- Event settings get min_rating_games: how many ended games a player needs
-- for a place in the rating. 1 lists everyone who has played, as before.
UPDATE events SET settings = settings || '{"min_rating_games": 1}';
