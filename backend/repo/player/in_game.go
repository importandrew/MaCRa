package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) IsInLiveGame(ctx context.Context, id int) (bool, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	inGameStatuses := []models.GameStatus{
		models.GameCreated,
		models.GameLive,
	}

	rows, err := conn.Query(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM games g
			JOIN game_players gp ON g.id = gp.game_id
			WHERE gp.player_id = $1
			AND g.status = ANY($2)
		)`,
		id, inGameStatuses)
	if err != nil {
		return false, err
	}

	exists, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return false, err
	}

	return exists, nil
}
