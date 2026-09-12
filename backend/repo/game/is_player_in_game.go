package game

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *Repo) IsPlayerInGame(ctx context.Context, gameID, playerID int) (bool, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM game_players
			WHERE game_id = $1 AND player_id = $2
		)`,
		gameID, playerID,
	)
	if err != nil {
		return false, err
	}

	return pgx.CollectOneRow(rows, pgx.RowTo[bool])
}
