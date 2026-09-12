package game

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) GameResult(ctx context.Context, gameID int) (models.GameResult, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return models.GameResult{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx,
		`SELECT finished_at FROM games WHERE id = $1`,
		gameID,
	)
	if err != nil {
		return models.GameResult{}, err
	}

	finishedAt, err := pgx.CollectOneRow(rows, pgx.RowTo[sql.NullTime])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.GameResult{}, models.ErrGameNotFound
		}
		return models.GameResult{}, err
	}

	rows, err = conn.Query(ctx, `
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

	if finishedAt.Valid {
		result.FinishedAt = &finishedAt.Time

		if result.FinishedAt.Add(models.UndoFinalRoundTimeout).After(time.Now()) {
			result.CanUndoLastRound = true
		}
	}

	return result, nil
}

type batchPlayerResultRow struct {
	GameID int     `db:"game_id"`
	Name   string  `db:"name"`
	ID     int     `db:"id"`
	TP     float32 `db:"tp"`
	GP     int     `db:"gp"`
}

type batchFinishedAtRow struct {
	GameID     int          `db:"id"`
	FinishedAt sql.NullTime `db:"finished_at"`
}

func (r *Repo) GameResults(ctx context.Context, gameIDs []int) (map[int]models.GameResult, error) {
	if len(gameIDs) == 0 {
		return map[int]models.GameResult{}, nil
	}

	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	finishedRows, err := conn.Query(ctx, `
		SELECT id, finished_at
		FROM games
		WHERE id = ANY($1)`,
		gameIDs,
	)
	if err != nil {
		return nil, err
	}

	finished, err := pgx.CollectRows(finishedRows, pgx.RowToStructByName[batchFinishedAtRow])
	if err != nil {
		return nil, err
	}

	results := make(map[int]models.GameResult, len(finished))
	for _, row := range finished {
		result := models.GameResult{
			ID:            row.GameID,
			PlayerResults: []models.PlayerResult{},
		}
		if row.FinishedAt.Valid {
			t := row.FinishedAt.Time
			result.FinishedAt = &t
		}
		results[row.GameID] = result
	}

	playerRows, err := conn.Query(ctx, `
		SELECT
			pr.game_id,
			pr.tournament_points AS tp,
			pr.game_points AS gp,
			u.name,
			u.id
		FROM player_results pr
		JOIN users u ON pr.player_id = u.id
		WHERE pr.game_id = ANY($1)
		ORDER BY pr.game_id, pr.tournament_points DESC`,
		gameIDs,
	)
	if err != nil {
		return nil, err
	}

	playerResults, err := pgx.CollectRows(playerRows, pgx.RowToStructByName[batchPlayerResultRow])
	if err != nil {
		return nil, err
	}

	for _, pr := range playerResults {
		result, ok := results[pr.GameID]
		if !ok {
			result = models.GameResult{
				ID:            pr.GameID,
				PlayerResults: []models.PlayerResult{},
			}
		}
		result.PlayerResults = append(result.PlayerResults, models.PlayerResult{
			ID:   pr.ID,
			Name: pr.Name,
			TP:   pr.TP,
			GP:   pr.GP,
		})
		results[pr.GameID] = result
	}

	return results, nil
}
