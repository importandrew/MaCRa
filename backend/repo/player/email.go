package player

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) IsEmailVerified(ctx context.Context, userID int) (bool, error) {
	rows, err := r.db.Query(ctx, `
		SELECT email_verified_at
		FROM users
		WHERE id = $1`,
		userID,
	)
	if err != nil {
		return false, err
	}

	verifiedAt, err := pgx.CollectOneRow(rows, pgx.RowTo[*time.Time])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, models.ErrPlayerNotFound
		}
		return false, err
	}

	return verifiedAt != nil, nil
}

func (r *Repo) SetEmailVerified(ctx context.Context, userID int) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE users
		SET email_verified_at = $1
		WHERE id = $2
		AND email_verified_at IS NULL`,
		time.Now().UTC(), userID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		rows, err := r.db.Query(ctx, `
			SELECT EXISTS(
				SELECT 1
				FROM users
				WHERE id = $1
			)`,
			userID,
		)
		if err != nil {
			return err
		}

		exists, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
		if err != nil {
			return err
		}
		if !exists {
			return models.ErrPlayerNotFound
		}
	}

	return nil
}

func (r *Repo) UpdatePassword(ctx context.Context, userID int, hash string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE users
		SET password = $1
		WHERE id = $2`,
		hash, userID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return models.ErrPlayerNotFound
	}

	return nil
}
