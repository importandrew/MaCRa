package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// RulesFrozen reports whether the event's game rules can no longer change:
// it has a game that wasn't aborted.
func (r *Repo) RulesFrozen(ctx context.Context, eventID int) (bool, error) {
	rows, err := r.db.Query(ctx, `
		SELECT EXISTS (SELECT 1 FROM games WHERE event_id = $1 AND status <> $2)`,
		eventID, models.GameAborted,
	)
	if err != nil {
		return false, err
	}

	return pgx.CollectOneRow(rows, pgx.RowTo[bool])
}

// UpdateRules sets the event's game rules. It returns
// models.ErrEventAlreadyFinished for a finished event and
// models.ErrEventRulesFrozen once the event has a game (see RulesFrozen).
//
// The event row is locked FOR UPDATE before the games are checked. Game
// creation locks the same row FOR SHARE and copies the rules in the same
// transaction, so a game can't start by rules that are being changed.
func (r *Repo) UpdateRules(ctx context.Context, eventID int, rules models.GameRules) error {
	return r.InTx(ctx, func(tx *Repo) error {
		rows, err := tx.db.Query(ctx, `SELECT finished_at FROM events WHERE id = $1 FOR UPDATE`, eventID)
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

		frozen, err := tx.RulesFrozen(ctx, eventID)
		if err != nil {
			return fmt.Errorf("check games: %w", err)
		}
		if frozen {
			return models.ErrEventRulesFrozen
		}

		if _, err := tx.db.Exec(ctx, `UPDATE events SET rules = $1 WHERE id = $2`, rules, eventID); err != nil {
			return fmt.Errorf("update rules: %w", err)
		}

		return nil
	})
}
