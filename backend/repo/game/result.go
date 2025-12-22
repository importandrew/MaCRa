package game

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) GameResult(ctx context.Context, gameID int) (models.GameResult, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return models.GameResult{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT 
			pr.tournament_points AS tp,
			pr.game_points AS gp,
			u.name,
			u.id
		FROM player_results pr
		JOIN users u ON pr.player_id = u.id
		WHERE pr.game_id = $1
		ORDER BY pr.tournament_points DESC`,
		gameID)
	if err != nil {
		return models.GameResult{}, err
	}

	playerResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.PlayerResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.GameResult{}, models.ErrGameNotFound
		}
		return models.GameResult{}, err
	}

	result := models.GameResult{
		ID:            gameID,
		PlayerResults: playerResults,
	}

	return result, nil
}
