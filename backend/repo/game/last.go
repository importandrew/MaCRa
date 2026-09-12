package game

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) LastGame(ctx context.Context, playerID, eventID int) (int, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT g.id 
		FROM games g
		JOIN game_players gp ON g.id = gp.game_id
		WHERE gp.player_id = $1 AND g.event_id = $2 AND g.status = $3
		ORDER BY g.finished_at DESC NULLS LAST
		FETCH FIRST 1 ROWS ONLY`,
		playerID, eventID, models.GameEnded)
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
