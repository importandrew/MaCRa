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
// that ended less than models.UndoFinalRoundTimeout ago; such a game becomes
// live again and loses its results. Everything happens in one transaction.
//
// logSize is the number of log entries the caller saw: as with
// SaveRoundOutcome, the game row is locked and the undo is refused with
// models.ErrLogChanged if the log has changed since, so two referees pressing
// Undo at once remove one entry, not two.
//
// Reviving an ended game takes the same advisory locks as NewGame and refuses
// with models.ErrUndoFinalRound if one of its players has started another
// live game in the event: a player can't have two live games.
func (r *Repo) UndoLastRound(ctx context.Context, gameID, logSize int) error {
	players, err := r.Players(ctx, gameID)
	if err != nil {
		return fmt.Errorf("players: %w", err)
	}

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

	if err := lockUnchangedLog(ctx, tx, gameID, logSize); err != nil {
		return err
	}

	type gameState struct {
		Status     models.GameStatus `db:"status"`
		FinishedAt *time.Time        `db:"finished_at"`
		EventID    int               `db:"event_id"`
	}

	rows, err := tx.Query(ctx, `SELECT status, finished_at, event_id FROM games WHERE id = $1`, gameID)
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
		if game.FinishedAt != nil && game.FinishedAt.Add(models.UndoFinalRoundTimeout).Before(time.Now()) {
			return models.ErrTooLate
		}
		seating := [4]int{players.East, players.South, players.West, players.North}
		if err := reviveGame(ctx, tx, gameID, game.EventID, seating); err != nil {
			return err
		}
	default:
		return fmt.Errorf("undo in a %s game: %w", game.Status, models.ErrGameNotLive)
	}

	if err := deleteLastRound(ctx, tx, gameID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// reviveGame makes an ended game live again: it deletes the players' results
// and clears finished_at. It refuses with models.ErrUndoFinalRound if one of
// the players already has another live game in the event, checked under the
// same advisory locks NewGame takes.
func reviveGame(ctx context.Context, tx pgx.Tx, gameID, eventID int, seating [4]int) error {
	if err := lockPlayersInEvent(ctx, tx, seating, eventID); err != nil {
		return err
	}

	rows, err := tx.Query(ctx, `
		SELECT u.name
		FROM games g
		JOIN game_players gp ON gp.game_id = g.id
		JOIN users u ON u.id = gp.player_id
		WHERE gp.player_id = ANY($1) AND g.event_id = $2 AND g.status = $3 AND g.id <> $4
		ORDER BY u.name
		FETCH FIRST 1 ROWS ONLY`,
		seating[:], eventID, models.GameLive, gameID)
	if err != nil {
		return fmt.Errorf("check live games: %w", err)
	}
	name, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
	switch {
	case err == nil:
		return fmt.Errorf("%w: %s has already started another game", models.ErrUndoFinalRound, name)
	case !errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("check live games: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM player_results WHERE game_id = $1`, gameID); err != nil {
		return fmt.Errorf("delete player results: %w", err)
	}

	if _, err := tx.Exec(ctx, `UPDATE games SET status = $1, finished_at = NULL WHERE id = $2`,
		models.GameLive, gameID); err != nil {
		return fmt.Errorf("set game live: %w", err)
	}

	return nil
}

// deleteLastRound deletes the game's last round log entry, or returns
// models.ErrNoRounds if the log is empty.
func deleteLastRound(ctx context.Context, tx pgx.Tx, gameID int) error {
	rows, err := tx.Query(ctx, `
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

	if _, err := tx.Exec(ctx, `DELETE FROM rounds WHERE id = $1`, roundID); err != nil {
		return fmt.Errorf("delete round: %w", err)
	}

	return nil
}
