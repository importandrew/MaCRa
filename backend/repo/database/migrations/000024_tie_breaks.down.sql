UPDATE events SET tournament = tournament - 'ties' - 'tie_breakers'
WHERE tournament IS NOT NULL;
