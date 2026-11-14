package game

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

func (r *Repo) CurrentGame(ctx context.Context, playerID, eventID int) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.id 
		FROM games g
		JOIN game_players gp ON g.id = gp.game_id
		WHERE gp.player_id = $1 AND g.event_id = $2 AND g.status = $3
		ORDER BY g.created_at DESC NULLS LAST
		LIMIT 1`,
		playerID, eventID, models.GameLive)
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

// SeatedGame returns the player's game in the event that isn't over: the
// live one, or a tournament game whose tour is published but not started
// (status created). Round entries use CurrentGame, live games only.
func (r *Repo) SeatedGame(ctx context.Context, playerID, eventID int) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT g.id
		FROM games g
		JOIN game_players gp ON g.id = gp.game_id
		WHERE gp.player_id = $1 AND g.event_id = $2 AND g.status IN ($3, $4)
		ORDER BY g.status = $3 DESC, g.created_at DESC NULLS LAST
		LIMIT 1`,
		playerID, eventID, models.GameLive, models.GameCreated)
	if err != nil {
		return 0, err
	}

	gameID, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, models.ErrGameNotFound
	}
	return gameID, err
}

// TourSeat returns the tour number and table of a tournament game, nil for
// any other game.
func (r *Repo) TourSeat(ctx context.Context, gameID int) (*models.TourSeat, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.number AS tour, g.table_number AS table
		FROM games g JOIN tours t ON t.id = g.tour_id
		WHERE g.id = $1`, gameID)
	if err != nil {
		return nil, err
	}

	seat, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[tourSeatRow])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //nolint:nilnil // no tour: not a tournament game
	}
	if err != nil {
		return nil, err
	}
	return &models.TourSeat{Tour: seat.Tour, Table: seat.Table}, nil
}

type tourSeatRow struct {
	Tour  int `db:"tour"`
	Table int `db:"table"`
}
