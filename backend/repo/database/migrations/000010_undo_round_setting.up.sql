-- Event settings get undo_round: who may undo the last round. Existing
-- events get referees only, as before.
UPDATE events SET settings = settings || '{"undo_round": "referees"}';
