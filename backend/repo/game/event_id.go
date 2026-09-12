package game

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) EventID(ctx context.Context, gameID int) (int, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT event_id FROM games WHERE id = $1`, gameID)
	if err != nil {
		return 0, err
	}

	eventID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, models.ErrGameNotFound
		}
		return 0, err
	}

	return eventID, nil
}
