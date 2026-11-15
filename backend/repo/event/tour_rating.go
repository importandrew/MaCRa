package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// TourResults returns the players' results in the tournament's finished
// tours (aborted ones don't count), by tour and table, and the number of the
// live tour (0 when none).
func (r *Repo) TourResults(ctx context.Context, eventID int) ([]models.TourResultRow, int, error) {
	var (
		results    []models.TourResultRow
		inProgress int
	)
	err := r.InTx(ctx, func(tx *Repo) error {
		rows, err := tx.db.Query(ctx, `
			SELECT t.number AS tour, pr.player_id, u.name, pr.tournament_points AS tp,
				pr.game_points AS gp, g.table_number
			FROM tours t
			JOIN games g ON g.tour_id = t.id AND g.status = $3
			JOIN player_results pr ON pr.game_id = g.id
			JOIN users u ON u.id = pr.player_id
			WHERE t.event_id = $1 AND t.status = $2
			ORDER BY t.number, g.table_number, pr.player_id`,
			eventID, models.TourFinished, models.GameEnded)
		if err != nil {
			return fmt.Errorf("tour results: %w", err)
		}
		if results, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.TourResultRow]); err != nil {
			return fmt.Errorf("tour results: %w", err)
		}

		rows, err = tx.db.Query(ctx, `
			SELECT COALESCE(max(number), 0) FROM tours WHERE event_id = $1 AND status = $2`,
			eventID, models.TourLive)
		if err != nil {
			return fmt.Errorf("live tour: %w", err)
		}
		inProgress, err = pgx.CollectOneRow(rows, pgx.RowTo[int])
		if err != nil {
			return fmt.Errorf("live tour: %w", err)
		}
		return nil
	})
	return results, inProgress, err
}

// LiveTour returns the tournament's live tour and whether the player's own
// game in it has ended, or models.ErrTourNotFound when no tour is live.
func (r *Repo) LiveTour(ctx context.Context, eventID, playerID int) (models.LiveTour, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.number, EXISTS (
			SELECT 1 FROM games g JOIN game_players gp ON gp.game_id = g.id
			WHERE g.tour_id = t.id AND g.status = $3 AND gp.player_id = $2) AS finished
		FROM tours t WHERE t.event_id = $1 AND t.status = $4
		LIMIT 1`,
		eventID, playerID, models.GameEnded, models.TourLive)
	if err != nil {
		return models.LiveTour{}, fmt.Errorf("live tour: %w", err)
	}
	tour, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.LiveTour])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.LiveTour{}, models.ErrTourNotFound
	case err != nil:
		return models.LiveTour{}, fmt.Errorf("live tour: %w", err)
	}
	return tour, nil
}
