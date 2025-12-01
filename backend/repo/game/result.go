package game

import (
	"context"
	"errors"
	"fmt"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

var errGameNotFound = errors.New("game not found")

type PlayerResult struct {
	TP   float32 `db:"tp"`
	GP   int     `db:"gp"`
	Name string  `db:"name"`
	ID   int     `db:"id"`
}

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
	defer rows.Close()

	playerResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[PlayerResult])
	if err != nil {
		return models.GameResult{}, err
	}

	if len(playerResults) == 0 {
		return models.GameResult{}, errGameNotFound
	}

	result := models.GameResult{
		ID:            gameID,
		PlayerResults: make([]models.PlayerResult, 0, len(playerResults)),
	}

	for _, res := range playerResults {
		result.PlayerResults = append(result.PlayerResults, models.PlayerResult{
			Name: res.Name,
			ID:   res.ID,
			TP:   res.TP,
			GP:   res.GP,
		})
	}

	return result, nil
}

func (r *Repo) LastGameResultByUserID(ctx context.Context, id int) (models.GameResult, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return models.GameResult{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT g.id 
		FROM games g
		JOIN game_players gp ON g.id = gp.game_id
		WHERE gp.player_id = $1 AND g.status = $2
		ORDER BY g.finished_at DESC
		FETCH FIRST 1 ROWS ONLY`,
		id, models.GameEnded)
	if err != nil {
		return models.GameResult{}, err
	}
	gameID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.GameResult{}, fmt.Errorf("no games for player with ID %d found", id)
		}
		return models.GameResult{}, err
	}

	finalResult, err := r.GameResult(ctx, gameID)
	if err != nil {
		return models.GameResult{}, err
	}

	return finalResult, nil
}
