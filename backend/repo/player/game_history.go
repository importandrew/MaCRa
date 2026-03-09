package player

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) GameHistory(ctx context.Context, playerID int) ([]int, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT g.id
		FROM games AS g
		JOIN game_players AS gp ON gp.game_id = g.id
		WHERE g.status = $1
  			AND gp.player_id = $2
		ORDER BY g.id DESC`,
		models.GameEnded, playerID,
	)
	if err != nil {
		return nil, err
	}

	gameIDs, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return gameIDs, nil
}
