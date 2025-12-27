package game

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

// SetStatus sets game status.
func (r *Repo) Status(ctx context.Context, gameID int) (models.GameStatus, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT status FROM games WHERE id = $1`, gameID)
	if err != nil {
		return "", err
	}

	status, err := pgx.CollectOneRow(rows, pgx.RowTo[models.GameStatus])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrGameNotFound
		}
		return "", err
	}

	return status, nil
}

// SetStatus sets game status.
func (r *Repo) SetStatus(ctx context.Context, gameID int, status models.GameStatus) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	res, err := conn.Exec(ctx, `UPDATE games SET status = $1 WHERE id = $2`, status, gameID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return models.ErrGameNotFound
	}

	return nil
}
