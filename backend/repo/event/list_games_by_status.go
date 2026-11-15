package event

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// ListGamesByStatus returns a page of the event's games with the status,
// newest first. hiddenTourID (0: none) leaves out the games of that tour
// except the ones callerID plays in: a live tour's tables are hidden from a
// player until their own game ends (docs/user-scenarios.md, area 8).
func (r *Repo) ListGamesByStatus(
	ctx context.Context,
	eventID int,
	status models.GameStatus,
	hiddenTourID, callerID int,
	limit, offset int,
) ([]models.EventGameItemRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.id, g.status
		FROM games g
		WHERE g.event_id = $1
			AND g.status = $2
			-- The tables of a hidden tour, except the caller's own.
			AND (g.tour_id IS DISTINCT FROM $5 OR EXISTS (
				SELECT 1 FROM game_players gp WHERE gp.game_id = g.id AND gp.player_id = $6))
		ORDER BY g.id DESC
		LIMIT $3 OFFSET $4`,
		eventID, status, limit, offset, hiddenTourID, callerID,
	)
	if err != nil {
		return nil, err
	}

	games, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.EventGameItemRow])
	if err != nil {
		return nil, err
	}

	return games, nil
}
