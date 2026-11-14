package event

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// List returns every event with the caller's roles in it (none for a
// visitor, userID 0) and the start of its last game that wasn't aborted, in
// one query. Active events come first, the ones played most recently (or,
// without games, created) on top; finished events follow, the latest
// finished on top.
func (r *Repo) List(ctx context.Context, userID int) ([]models.EventListRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT e.id, e.type, e.name, e.description, e.owner_user_id,
			e.visibility, e.ratings_reach, e.games_reach, e.rules, e.settings,
			e.created_at, e.finished_at,
			COALESCE(
				(SELECT array_agg(eur.role ORDER BY eur.role) FROM event_user_roles eur
				 WHERE eur.event_id = e.id AND eur.user_id = $1),
				'{}') AS roles,
			(SELECT max(g.created_at) FROM games g
			 WHERE g.event_id = e.id AND g.status <> $2) AS last_game_at
		FROM events e
		ORDER BY e.finished_at IS NOT NULL,
			COALESCE(e.finished_at,
				(SELECT max(g.created_at) FROM games g WHERE g.event_id = e.id AND g.status <> $2),
				e.created_at) DESC,
			e.id DESC`,
		userID, models.GameAborted,
	)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.EventListRow])
}
