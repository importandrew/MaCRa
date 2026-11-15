package game

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) LastGame(ctx context.Context, playerID, eventID int) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.id
		FROM games g
		WHERE g.event_id = $2 AND g.status = $3
			-- Their own game, or one they sat at as a substitute.
			AND (EXISTS (SELECT 1 FROM game_players gp WHERE gp.game_id = g.id AND gp.player_id = $1)
				OR EXISTS (SELECT 1 FROM game_seats s WHERE s.game_id = g.id AND s.sitter_id = $1))
		ORDER BY g.finished_at DESC NULLS LAST
		LIMIT 1`,
		playerID, eventID, models.GameEnded)
	if err != nil {
		return 0, err
	}

	gameID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, models.ErrGameNotFound
		}
		return 0, err
	}

	return gameID, nil
}
