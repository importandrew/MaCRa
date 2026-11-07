package game

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) SaveRoundOutcome(ctx context.Context, gameID int, data models.RoundData) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	players, err := r.Players(ctx, gameID)
	if err != nil {
		return fmt.Errorf("players: %w", err)
	}

	// Start transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := lockUnchangedLog(ctx, tx, gameID, data.LogSize); err != nil {
		return err
	}

	// Insert round event (event_sequence is auto-incremented by database trigger)
	type roundInsertResult struct {
		ID            int `db:"id"`
		EventSequence int `db:"event_sequence"`
	}

	// Convert zero values to NULL for foreign key constraints
	winnerID := sql.NullInt64{Valid: data.WinnerID != 0, Int64: int64(data.WinnerID)}
	loserID := sql.NullInt64{Valid: data.LoserID != 0, Int64: int64(data.LoserID)}
	value := sql.NullInt64{Valid: data.Value != 0, Int64: int64(data.Value)}

	rows, err := tx.Query(ctx, `
		INSERT INTO rounds(
			game_id, round_number, event_sequence, outcome, winner_id, loser_id, 
			value, force_next_round, created_at)
		VALUES($1, $2, NULL, $3, $4, $5, $6, $7, $8)
		RETURNING id, event_sequence`,
		gameID, data.Round, data.Outcome, winnerID, loserID,
		value, data.ForceNextRound, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("insert round: %w", err)
	}
	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[roundInsertResult])
	if err != nil {
		return fmt.Errorf("insert round: %w", err)
	}
	roundID := result.ID
	eventSequence := result.EventSequence

	// Prepare values for batch insert
	args := pgx.NamedArgs{
		"round_id":  roundID,
		"game_id":   gameID,
		"round_num": data.Round,
		"event_seq": eventSequence,
		"e_score":   data.EastStartScore,
		"e_id":      players.East,
		"s_score":   data.SouthStartScore,
		"s_id":      players.South,
		"w_score":   data.WestStartScore,
		"w_id":      players.West,
		"n_score":   data.NorthStartScore,
		"n_id":      players.North,
	}

	// Batch insert all round scores at once
	_, err = tx.Exec(ctx, `
		INSERT INTO round_scores(round_id, game_id, round_number, event_sequence, player_id, score)
		VALUES
			(@round_id, @game_id, @round_num, @event_seq, @e_id, @e_score),
			(@round_id, @game_id, @round_num, @event_seq, @s_id, @s_score),
			(@round_id, @game_id, @round_num, @event_seq, @w_id, @w_score),
			(@round_id, @game_id, @round_num, @event_seq, @n_id, @n_score)
		ON CONFLICT (round_id, player_id) DO UPDATE SET score = EXCLUDED.score`,
		args)
	if err != nil {
		return fmt.Errorf("batch insert round scores: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// lockUnchangedLog locks the game row, so that concurrent entries to one game
// are saved one at a time, and refuses the entry with ErrLogChanged
// unless the game still has logSize round log entries. The entry was computed
// from a game replayed from that many entries; saving it after another entry
// got in first would apply it to a game it wasn't computed for (a
// double-tapped penalty would be fined twice).
func lockUnchangedLog(ctx context.Context, tx pgx.Tx, gameID, logSize int) error {
	if _, err := tx.Exec(ctx, `SELECT 1 FROM games WHERE id = $1 FOR UPDATE`, gameID); err != nil {
		return fmt.Errorf("lock game: %w", err)
	}

	rows, err := tx.Query(ctx, `SELECT count(*) FROM rounds WHERE game_id = $1`, gameID)
	if err != nil {
		return fmt.Errorf("count rounds: %w", err)
	}
	current, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return fmt.Errorf("count rounds: %w", err)
	}

	if current != logSize {
		return fmt.Errorf("%w (computed from %d log entries, the game has %d)", models.ErrLogChanged, logSize, current)
	}

	return nil
}
