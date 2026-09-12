package game

import (
	"context"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) NewGame(ctx context.Context, seating [4]int, eventID int) (int, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	east := seating[0]
	south := seating[1]
	west := seating[2]
	north := seating[3]

	// Start transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	now := time.Now().UTC()

	// Insert game and get ID
	rows, err := tx.Query(ctx, `
		INSERT INTO games(created_at, status, event_id) 
		VALUES($1, $2, $3)
		RETURNING id`,
		now, models.GameLive, eventID)
	if err != nil {
		return 0, err
	}
	gameID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, err
	}

	// Insert game players
	_, err = tx.Exec(ctx, `
		INSERT INTO game_players(game_id, player_id, starting_wind) 
		VALUES 
			($1, $2, 'east'),
			($1, $3, 'south'),
			($1, $4, 'west'),
			($1, $5, 'north')`,
		gameID, east, south, west, north)
	if err != nil {
		return 0, err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return gameID, nil
}
