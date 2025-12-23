package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

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
