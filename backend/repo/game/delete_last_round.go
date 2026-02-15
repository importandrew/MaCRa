package game

import (
	"context"
	"errors"
	"fmt"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) DeleteLastRound(ctx context.Context, gameID int) error {
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

	rows, err := tx.Query(ctx, `
		SELECT id FROM rounds
		WHERE game_id = $1
		ORDER BY round_number DESC, event_sequence DESC
		LIMIT 1`,
		gameID,
	)
	if err != nil {
		return err
	}

	roundID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrNoRounds
		}
		return fmt.Errorf("find last round: %w", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM rounds WHERE id = $1`, roundID)
	if err != nil {
		return fmt.Errorf("delete round: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
