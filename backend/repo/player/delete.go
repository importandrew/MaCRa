package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) Delete(ctx context.Context, email string) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	res, err := conn.Exec(ctx, "DELETE FROM users WHERE email = $1", email)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return models.ErrPlayerNotFound
	}

	return nil
}
