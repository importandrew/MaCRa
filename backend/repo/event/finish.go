package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// Finish marks the event finished. It returns models.ErrEventAlreadyFinished
// if it already is, and models.ErrEventHasOngoingGames if a game in it is
// live.
//
// The event row is locked FOR UPDATE before the live games are checked. Game
// creation and the undo that revives an ended game lock the same row FOR
// SHARE, so a game can't become live while the event is being finished.
func (r *Repo) Finish(ctx context.Context, eventID int) error {
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

	rows, err := tx.Query(ctx, `SELECT finished_at FROM events WHERE id = $1 FOR UPDATE`, eventID)
	if err != nil {
		return fmt.Errorf("lock event: %w", err)
	}
	finishedAt, err := pgx.CollectOneRow(rows, pgx.RowTo[*time.Time])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.ErrEventNotFound
	case err != nil:
		return fmt.Errorf("lock event: %w", err)
	case finishedAt != nil:
		return models.ErrEventAlreadyFinished
	}

	rows, err = tx.Query(ctx, `
		SELECT EXISTS (SELECT 1 FROM games WHERE event_id = $1 AND status = $2)`,
		eventID, models.GameLive,
	)
	if err != nil {
		return fmt.Errorf("check live games: %w", err)
	}
	hasLive, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return fmt.Errorf("check live games: %w", err)
	}
	if hasLive {
		return models.ErrEventHasOngoingGames
	}

	if _, err := tx.Exec(ctx, `UPDATE events SET finished_at = $1 WHERE id = $2`,
		time.Now().UTC(), eventID); err != nil {
		return fmt.Errorf("finish event: %w", err)
	}

	return tx.Commit(ctx)
}
