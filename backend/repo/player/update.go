package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) Update(ctx context.Context, id int, name string) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	res, err := conn.Exec(ctx, "UPDATE users SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return models.ErrPlayerNotFound
	}

	return nil
}
