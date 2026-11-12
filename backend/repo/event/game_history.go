package event

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) GameHistory(ctx context.Context, eventID int) ([]int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id
		FROM games
		WHERE status = $1
  			AND event_id = $2
		ORDER BY id DESC`,
		models.GameEnded, eventID,
	)
	if err != nil {
		return nil, err
	}

	gameIDs, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return nil, err
	}

	return gameIDs, nil
}

func (r *Repo) PlayerGameHistory(ctx context.Context, eventID, playerID int) ([]int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.id
		FROM games AS g
		JOIN game_players AS gp ON gp.game_id = g.id
		WHERE g.status = $1
  			AND gp.player_id = $2
  			AND g.event_id = $3
		ORDER BY g.id DESC`,
		models.GameEnded, playerID, eventID,
	)
	if err != nil {
		return nil, err
	}

	gameIDs, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return nil, err
	}

	return gameIDs, nil
}
