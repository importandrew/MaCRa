package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListEndedPlayerResults(ctx context.Context, eventID int) ([]models.EventPlayerResultRow, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT
			pr.player_id,
			u.name,
			pr.tournament_points AS tp,
			pr.game_points AS gp
		FROM player_results pr
		JOIN games g ON g.id = pr.game_id
		JOIN users u ON u.id = pr.player_id
		WHERE g.event_id = $1
			AND g.status = $2
		ORDER BY pr.player_id, pr.game_id`,
		eventID, models.GameEnded,
	)
	if err != nil {
		return nil, err
	}

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.EventPlayerResultRow])
	if err != nil {
		return nil, err
	}

	return results, nil
}
