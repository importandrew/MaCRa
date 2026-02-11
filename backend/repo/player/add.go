package player

import (
	"context"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) Add(ctx context.Context, email, name, hash string) (int, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		"INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id",
		name,
		hash,
		email,
	)
	if err != nil {
		return 0, err
	}

	playerID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		if repo.IsTargetErr(23505, err) {
			return 0, models.ErrPlayerExists
		}

		if repo.IsTargetErr(22001, err) {
			return 0, models.ErrLongName
		}

		return 0, err
	}

	return playerID, nil
}
