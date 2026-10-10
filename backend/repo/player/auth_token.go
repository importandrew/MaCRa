package player

import (
	"context"
	"errors"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) InvalidateAuthTokens(ctx context.Context, userID int, purpose models.AuthPurpose) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
		UPDATE auth_tokens
			SET used_at = $1
		WHERE user_id = $2
			AND purpose = $3
			AND used_at IS NULL`,
		time.Now().UTC(), userID, purpose,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) CreateAuthToken(ctx context.Context, userID int, purpose models.AuthPurpose, tokenHash string, expiresAt time.Time) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
		INSERT INTO auth_tokens (user_id, purpose, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		userID, purpose, tokenHash, expiresAt, time.Now().UTC(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetAuthTokenByHash(ctx context.Context, tokenHash string) (models.AuthToken, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return models.AuthToken{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT id, user_id, purpose, token_hash, expires_at, used_at, created_at
		FROM auth_tokens
		WHERE token_hash = $1`,
		tokenHash,
	)
	if err != nil {
		return models.AuthToken{}, err
	}

	token, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.AuthToken])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.AuthToken{}, models.ErrInvalidAuthToken
		}
		return models.AuthToken{}, err
	}

	return token, nil
}

func (r *Repo) MarkAuthTokenUsed(ctx context.Context, id int) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tag, err := conn.Exec(ctx, `
		UPDATE auth_tokens
		SET used_at = $1
		WHERE id = $2
		AND used_at IS NULL`,
		time.Now().UTC(), id,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return models.ErrInvalidAuthToken
	}

	return nil
}
