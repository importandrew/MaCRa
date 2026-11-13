package game

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// Summary returns what the capability rules need to know about a game: its
// status, when it ended and how many round log entries it has. It returns
// models.ErrGameNotFound for an unknown game.
func (r *Repo) Summary(ctx context.Context, gameID int) (models.GameSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			g.status,
			g.finished_at,
			(SELECT count(*) FROM rounds WHERE game_id = g.id) AS log_size
		FROM games g
		WHERE g.id = $1`,
		gameID)
	if err != nil {
		return models.GameSummary{}, err
	}

	summary, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.GameSummary])
	if errors.Is(err, pgx.ErrNoRows) {
		return models.GameSummary{}, models.ErrGameNotFound
	}
	return summary, err
}

// Abort makes a live game aborted. With onlyWithoutRounds it does so only
// while the game has no round log entries, in the same statement, so a round
// entered at the same moment can't slip in between the check and the abort;
// otherwise it returns models.ErrGameHasRounds.
func (r *Repo) Abort(ctx context.Context, gameID int, onlyWithoutRounds bool) error {
	res, err := r.db.Exec(ctx, `
		UPDATE games SET status = $2
		WHERE id = $1 AND status = $3
			AND (NOT $4 OR NOT EXISTS (SELECT 1 FROM rounds WHERE game_id = $1))`,
		gameID, models.GameAborted, models.GameLive, onlyWithoutRounds)
	if err != nil {
		return err
	}
	if res.RowsAffected() > 0 {
		return nil
	}

	// Nothing changed: say why.
	status, err := r.Status(ctx, gameID)
	if err != nil {
		return err
	}
	if status != models.GameLive {
		return models.ErrGameNotLive
	}
	return models.ErrGameHasRounds
}
