package event

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) GameHistory(ctx context.Context, eventID int) ([]int, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return gameIDs, nil
}

func (r *Repo) PlayerGameHistory(ctx context.Context, eventID, playerID int) ([]int, error) {
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
  			AND g.event_id = $3
		ORDER BY g.id DESC`,
		models.GameEnded, playerID, eventID,
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
