package player

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) IsEmailVerified(ctx context.Context, userID int) (bool, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT email_verified_at
		FROM users
		WHERE id = $1`,
		userID,
	)
	if err != nil {
		return false, err
	}

	verifiedAt, err := pgx.CollectOneRow(rows, pgx.RowTo[sql.NullTime])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, models.ErrPlayerNotFound
		}
		return false, err
	}

	return verifiedAt.Valid, nil
}

func (r *Repo) SetEmailVerified(ctx context.Context, userID int) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tag, err := conn.Exec(ctx, `
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
		rows, err := conn.Query(ctx, `
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
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tag, err := conn.Exec(ctx, `
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
