package game

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// UndoLastRound deletes the last round log entry of a live game, or of a game
// that ended after endedAfter; such a game becomes live again and loses its
// results. An ended game that ended before endedAfter is refused with
// models.ErrTooLate. Everything happens in one transaction.
//
// logSize is the number of log entries the caller saw: as with
// SaveRoundOutcome, the game row is locked and the undo is refused with
// models.ErrLogChanged if the log has changed since, so two referees pressing
// Undo at once remove one entry, not two.
//
// Reviving an ended game takes the same advisory locks as NewGame and refuses
// with models.ErrPlayerInAnotherGame if one of its players has started
// another live game in the event: a player can't have two live games.
func (r *Repo) UndoLastRound(ctx context.Context, gameID, logSize int, endedAfter time.Time) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if err := tx.lockUnchangedLog(ctx, gameID, logSize); err != nil {
			return err
		}

		type gameState struct {
			Status     models.GameStatus `db:"status"`
			FinishedAt *time.Time        `db:"finished_at"`
			EventID    int               `db:"event_id"`
		}

		rows, err := tx.db.Query(ctx, `SELECT status, finished_at, event_id FROM games WHERE id = $1`, gameID)
		if err != nil {
			return fmt.Errorf("get game: %w", err)
		}
		game, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[gameState])
		if err != nil {
			return fmt.Errorf("get game: %w", err)
		}

		switch game.Status {
		case models.GameLive:
		case models.GameEnded:
			if game.FinishedAt != nil && game.FinishedAt.Before(endedAfter) {
				return models.ErrTooLate
			}
			players, err := tx.Players(ctx, gameID)
			if err != nil {
				return fmt.Errorf("players: %w", err)
			}
			seating := [4]int{players.East, players.South, players.West, players.North}
			if err := tx.reviveGame(ctx, gameID, game.EventID, seating); err != nil {
				return err
			}
		default:
			return fmt.Errorf("undo in a %s game: %w", game.Status, models.ErrGameNotLive)
		}

		return tx.deleteLastRound(ctx, gameID)
	})
}

// reviveGame makes an ended game live again: it deletes the players' results
// and clears finished_at. Like NewGame, it refuses with
// models.ErrEventAlreadyFinished in a finished event, and with
// models.ErrPlayerInAnotherGame if one of the players already has another live
// game in the event, checked under the same locks NewGame takes. Call it on a Repo
// from InTx.
func (r *Repo) reviveGame(ctx context.Context, gameID, eventID int, seating [4]int) error {
	if err := r.lockOpenEvent(ctx, eventID); err != nil {
		return err
	}

	if err := r.lockPlayersInEvent(ctx, seating, eventID); err != nil {
		return err
	}

	if err := r.checkNoOtherLiveGame(ctx, seating, eventID, gameID); err != nil {
		return err
	}

	if _, err := r.db.Exec(ctx, `DELETE FROM player_results WHERE game_id = $1`, gameID); err != nil {
		return fmt.Errorf("delete player results: %w", err)
	}

	if _, err := r.db.Exec(ctx, `UPDATE games SET status = $1, finished_at = NULL WHERE id = $2`,
		models.GameLive, gameID); err != nil {
		return fmt.Errorf("set game live: %w", err)
	}

	return nil
}

// deleteLastRound deletes the game's last round log entry, or returns
// models.ErrNoRounds if the log is empty.
func (r *Repo) deleteLastRound(ctx context.Context, gameID int) error {
	rows, err := r.db.Query(ctx, `
		SELECT id FROM rounds
		WHERE game_id = $1
		ORDER BY round_number DESC, event_sequence DESC
		LIMIT 1`,
		gameID,
	)
	if err != nil {
		return fmt.Errorf("find last round: %w", err)
	}

	roundID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrNoRounds
		}
		return fmt.Errorf("find last round: %w", err)
	}

	if _, err := r.db.Exec(ctx, `DELETE FROM rounds WHERE id = $1`, roundID); err != nil {
		return fmt.Errorf("delete round: %w", err)
	}

	return nil
}
