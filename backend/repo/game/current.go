package game

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) CurrentGame(ctx context.Context, playerID, eventID int) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.id 
		FROM games g
		JOIN game_players gp ON g.id = gp.game_id
		WHERE gp.player_id = $1 AND g.event_id = $2 AND g.status = $3
		ORDER BY g.created_at DESC NULLS LAST
		LIMIT 1`,
		playerID, eventID, models.GameLive)
	if err != nil {
		return 0, err
	}

	gameID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, models.ErrGameNotFound
		}
		return 0, err
	}

	return gameID, nil
}
