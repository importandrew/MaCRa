package game

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// IsPlayerInGame reports whether the user sits at the game now: at their
// own seat, or at one they took as a substitute (game_sitters). A player
// whose seat a substitute took doesn't.
func (r *Repo) IsPlayerInGame(ctx context.Context, gameID, playerID int) (bool, error) {
	rows, err := r.db.Query(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM game_sitters
			WHERE game_id = $1 AND sitter_id = $2
		)`,
		gameID, playerID,
	)
	if err != nil {
		return false, err
	}

	return pgx.CollectOneRow(rows, pgx.RowTo[bool])
}
