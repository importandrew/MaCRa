package player

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) GetPassByEmail(ctx context.Context, email string) (string, error) {
	rows, err := r.db.Query(ctx, `SELECT password FROM users WHERE email = $1`, email)
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
