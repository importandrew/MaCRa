-- A tournament's settings get what decides between players equal in TP and
-- GP (docs/user-scenarios.md, area 8, "Rating"): "share" the place (what
-- the rating did so far) or "break" the tie by tie_breakers in order, then
-- the draw number. The stored JSON stays complete.
UPDATE events SET tournament = tournament || '{"ties": "share", "tie_breakers": []}'::jsonb
WHERE tournament IS NOT NULL;
