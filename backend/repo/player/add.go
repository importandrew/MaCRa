package player

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo"
)

func (r *Repo) Add(ctx context.Context, email, name, hash string) (int, error) {
	rows, err := r.db.Query(
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
		if repo.IsTargetErr(repo.SQLStateUniqueViolation, err) {
			return 0, models.ErrPlayerExists
		}

		if repo.IsTargetErr(repo.SQLStateStringTooLong, err) {
			return 0, models.ErrLongName
		}

		return 0, err
	}

	return playerID, nil
}
