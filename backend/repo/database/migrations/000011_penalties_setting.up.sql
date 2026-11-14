-- Event settings get penalties: who may enter penalties and false mahjongs.
-- Existing events get the seated players too, as before.
UPDATE events SET settings = settings || '{"penalties": "seated"}';
