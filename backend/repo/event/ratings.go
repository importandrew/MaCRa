package event

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) ListEndedPlayerResults(ctx context.Context, eventID int) ([]models.EventPlayerResultRow, error) {
	rows, err := r.db.Query(ctx, `
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
			-- A tournament game counts once its tour is finished.
			AND (g.tour_id IS NULL OR EXISTS (
				SELECT 1 FROM tours t WHERE t.id = g.tour_id AND t.status = $3))
			-- A forfeited seat counts for nobody.
			AND NOT EXISTS (SELECT 1 FROM game_players f WHERE f.game_id = pr.game_id AND f.player_id = pr.player_id AND f.forfeit)
		ORDER BY pr.player_id, pr.game_id`,
		eventID, models.GameEnded, models.TourFinished,
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
