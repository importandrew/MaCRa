package game

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// UndoLastRound deletes the last round log entry of a live game, or of a game
// that ended after endedAfter; such a game becomes live again and loses its
// results. An ended game that ended before endedAfter is refused with
// models.ErrTooLate; with endedAfter nil every ended game is refused with
// models.ErrGameNotLive (a seated player's undo: the game may have ended
// since their check). Everything happens in one transaction.
//
// logSize is the number of log entries the caller saw: as with
// SaveRoundOutcome, the game row is locked and the undo is refused with
// models.ErrLogChanged if the log has changed since, so two referees pressing
// Undo at once remove one entry, not two.
//
// The deleted entry is kept in undone_rounds with undoneBy (the user undoing
// it) and the time, for the game's log.
//
// Reviving an ended game takes the same advisory locks as NewGame and refuses
// with models.ErrPlayerInAnotherGame if one of its players has started
// another live game in the event: a player can't have two live games.
func (r *Repo) UndoLastRound(ctx context.Context, gameID, logSize, undoneBy int, endedAfter *time.Time) error {
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
			if err := tx.reviveEndedGame(ctx, gameID, game.EventID, game.FinishedAt, endedAfter); err != nil {
				return err
			}
		default:
			return fmt.Errorf("undo in a %s game: %w", game.Status, models.ErrGameNotLive)
		}

		return tx.deleteLastRound(ctx, gameID, undoneBy)
	})
}

// reviveEndedGame makes the ended game live again for UndoLastRound, if it
// ended after endedAfter (models.ErrTooLate otherwise); with endedAfter nil
// it refuses with models.ErrGameNotLive. Call it on a Repo from InTx.
func (r *Repo) reviveEndedGame(ctx context.Context, gameID, eventID int, finishedAt, endedAfter *time.Time) error {
	if endedAfter == nil {
		return fmt.Errorf("undo in an ended game without reviving: %w", models.ErrGameNotLive)
	}
	if finishedAt != nil && finishedAt.Before(*endedAfter) {
		return models.ErrTooLate
	}

	players, err := r.Players(ctx, gameID)
	if err != nil {
		return fmt.Errorf("players: %w", err)
	}
	seating := [4]int{players.East, players.South, players.West, players.North}

	return r.reviveGame(ctx, gameID, eventID, seating)
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

	if err := r.checkPlayerRoles(ctx, seating, eventID); err != nil {
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

// deleteLastRound moves the game's last round log entry to undone_rounds,
// undone by undoneBy now, or returns models.ErrNoRounds if the log is empty.
func (r *Repo) deleteLastRound(ctx context.Context, gameID, undoneBy int) error {
	// The last entry and the one before it, whose scores are the last
	// entry's scores before (none: all zero).
	rows, err := r.db.Query(ctx, `
		SELECT id FROM rounds
		WHERE game_id = $1
		ORDER BY round_number DESC, event_sequence DESC
		LIMIT 2`,
		gameID,
	)
	if err != nil {
		return fmt.Errorf("find last round: %w", err)
	}

	ids, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return fmt.Errorf("find last round: %w", err)
	}
	if len(ids) == 0 {
		return models.ErrNoRounds
	}
	roundID := ids[0]
	var previousID *int
	if len(ids) > 1 {
		previousID = &ids[1]
	}

	// round_scores rows go with the round (ON DELETE CASCADE); the scores
	// subqueries still see them: every part of one statement reads the same
	// snapshot, taken before the DELETE.
	if _, err := r.db.Exec(ctx, `
		WITH gone AS (
			DELETE FROM rounds WHERE id = $1
			RETURNING game_id, round_number, outcome, winner_id, loser_id, value,
				force_next_round, entered_by, created_at
		)
		INSERT INTO undone_rounds (
			game_id, round_number, outcome, winner_id, loser_id, value,
			force_next_round, entered_by, entered_at, undone_by, undone_at,
			scores_before, scores_after)
		SELECT game_id, round_number, outcome, winner_id, loser_id, value,
			force_next_round, entered_by, created_at, $2, $3,
			COALESCE(`+seatScores("$4")+`, '{0,0,0,0}'), `+seatScores("$1")+`
		FROM gone`,
		roundID, nullIfZero(undoneBy), time.Now().UTC(), previousID,
	); err != nil {
		return fmt.Errorf("delete round: %w", err)
	}

	return nil
}

// seatScores is an SQL subquery for the players' scores after the round log
// entry whose ID is the parameter roundIDParam (e.g. "$1"), as an array by
// starting seat: east, south, west, north. NULL when there's no such entry.
func seatScores(roundIDParam string) string {
	return `(SELECT array_agg(rs.score ORDER BY CASE gp.starting_wind
			WHEN 'east' THEN 1 WHEN 'south' THEN 2 WHEN 'west' THEN 3 ELSE 4 END)
		FROM round_scores rs
		JOIN game_players gp ON gp.game_id = rs.game_id AND gp.player_id = rs.player_id
		WHERE rs.round_id = ` + roundIDParam + `::int)`
}
