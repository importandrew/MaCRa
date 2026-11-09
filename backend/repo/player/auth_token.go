package player

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
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

// IssueAuthToken stores a new token for the user and purpose and marks the
// unused ones of that purpose used. It returns models.ErrAuthMailLimit, and
// stores nothing, if the user already has as many tokens as a window of limits
// allows. The user row is locked FOR UPDATE, so concurrent requests are
// counted one after another.
func (r *Repo) IssueAuthToken(
	ctx context.Context,
	userID int,
	purpose models.AuthPurpose,
	tokenHash string,
	expiresAt time.Time,
	limits []models.AuthMailWindow,
) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `SELECT id FROM users WHERE id = $1 FOR UPDATE`, userID)
	if err != nil {
		return fmt.Errorf("lock user: %w", err)
	}
	if _, err := pgx.CollectOneRow(rows, pgx.RowTo[int]); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrPlayerNotFound
		}
		return fmt.Errorf("lock user: %w", err)
	}

	now := time.Now().UTC()
	for _, window := range limits {
		rows, err := tx.Query(ctx, `
			SELECT count(*) FROM auth_tokens
			WHERE user_id = $1 AND created_at > $2`,
			userID, now.Add(-window.Period),
		)
		if err != nil {
			return fmt.Errorf("count tokens: %w", err)
		}
		count, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
		if err != nil {
			return fmt.Errorf("count tokens: %w", err)
		}
		if count >= window.Max {
			return fmt.Errorf("%d in %s: %w", count, window.Period, models.ErrAuthMailLimit)
		}
	}

	if _, err := tx.Exec(ctx, `
		UPDATE auth_tokens SET used_at = $1
		WHERE user_id = $2 AND purpose = $3 AND used_at IS NULL`,
		now, userID, purpose,
	); err != nil {
		return fmt.Errorf("invalidate tokens: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO auth_tokens (user_id, purpose, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		userID, purpose, tokenHash, expiresAt, now,
	); err != nil {
		return fmt.Errorf("create token: %w", err)
	}

	return tx.Commit(ctx)
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
