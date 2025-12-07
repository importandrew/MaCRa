package player

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetPassByEmail(ctx context.Context, email string) (string, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT password FROM users WHERE email = $1`, email)
	if err != nil {
		return "", err
	}

	password, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrPlayerNotFound
		}
		return "", err
	}

	return password, nil
}
