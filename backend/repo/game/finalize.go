package game

import (
	"context"
	"fmt"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

// Finalize stores players' scores and sets game status to finished.
func (r *Repo) Finalize(ctx context.Context, data models.FinalGameData) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// Start transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	now := time.Now().UTC()

	// Prepare values for batch insert
	args := pgx.NamedArgs{
		"game_id":    data.GameID,
		"created_at": now,
		"e_id":       data.East.ID,
		"e_gp":       data.East.GP,
		"e_tp":       data.East.TP,
		"s_id":       data.South.ID,
		"s_gp":       data.South.GP,
		"s_tp":       data.South.TP,
		"w_id":       data.West.ID,
		"w_gp":       data.West.GP,
		"w_tp":       data.West.TP,
		"n_id":       data.North.ID,
		"n_gp":       data.North.GP,
		"n_tp":       data.North.TP,
	}

	// Batch insert all player results at once
	_, err = tx.Exec(ctx, `
		INSERT INTO player_results(game_id, player_id, game_points, tournament_points, created_at)
		VALUES
			(@game_id, @e_id, @e_gp, @e_tp, @created_at),
			(@game_id, @s_id, @s_gp, @s_tp, @created_at),
			(@game_id, @w_id, @w_gp, @w_tp, @created_at),
			(@game_id, @n_id, @n_gp, @n_tp, @created_at)
		ON CONFLICT (game_id, player_id) DO UPDATE SET
			game_points = EXCLUDED.game_points,
			tournament_points = EXCLUDED.tournament_points,
			created_at = EXCLUDED.created_at`,
		args)
	if err != nil {
		return fmt.Errorf("batch insert player results: %w", err)
	}

	// Update game status
	err = r.SetStatus(ctx, data.GameID, models.GameEnded)
	if err != nil {
		return fmt.Errorf("set status: %w", err)
	}

	// Set finished_at
	_, err = conn.Exec(ctx, `UPDATE games SET finished_at = $1 WHERE id = $2`, time.Now().UTC(), data.GameID)
	if err != nil {
		return fmt.Errorf("set finished_at for game %d: %w", data.GameID, err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
