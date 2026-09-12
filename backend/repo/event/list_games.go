package event

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) ListGames(ctx context.Context, eventID int) ([]models.Game, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT id, status, created_at, finished_at, event_id
		FROM games
		WHERE event_id = $1
		ORDER BY id`,
		eventID,
	)
	if err != nil {
		return nil, err
	}

	games, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Game])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return games, nil
}
