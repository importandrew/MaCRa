package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// SeatingPlan returns the tournament's seating prepared for all tours, or
// models.ErrNoSeatingPlan.
func (r *Repo) SeatingPlan(ctx context.Context, eventID int) (models.SeatingPlan, error) {
	rows, err := r.db.Query(ctx, `
		SELECT source, created_at, tours FROM seating_plans WHERE event_id = $1`, eventID)
	if err != nil {
		return models.SeatingPlan{}, fmt.Errorf("seating plan: %w", err)
	}
	plan, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.SeatingPlan])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.SeatingPlan{}, models.ErrNoSeatingPlan
	case err != nil:
		return models.SeatingPlan{}, fmt.Errorf("seating plan: %w", err)
	}
	return plan, nil
}

// SaveSeatingPlan stores the tournament's seating prepared for all tours,
// replacing an earlier one. It refuses a regular or finished event.
func (r *Repo) SaveSeatingPlan(ctx context.Context, eventID int, source models.SeatingPlanSource, tours models.PlanTours) error {
	return r.InTx(ctx, func(tx *Repo) error {
		isTournament, _, err := tx.tournamentState(ctx, eventID, false)
		switch {
		case err != nil:
			return err
		case !isTournament:
			return models.ErrNotTournament
		}

		_, err = tx.db.Exec(ctx, `
			INSERT INTO seating_plans (event_id, source, created_at, tours)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (event_id) DO UPDATE
			SET source = EXCLUDED.source, created_at = EXCLUDED.created_at, tours = EXCLUDED.tours`,
			eventID, source, time.Now().UTC(), tours)
		if err != nil {
			return fmt.Errorf("save seating plan: %w", err)
		}
		return nil
	})
}

// DeleteSeatingPlan deletes the tournament's seating prepared for all tours
// (the tours seated from it stay). It refuses a finished event and when
// there is none (models.ErrNoSeatingPlan).
func (r *Repo) DeleteSeatingPlan(ctx context.Context, eventID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if _, _, err := tx.tournamentState(ctx, eventID, false); err != nil {
			return err
		}
		res, err := tx.db.Exec(ctx, `DELETE FROM seating_plans WHERE event_id = $1`, eventID)
		if err != nil {
			return fmt.Errorf("delete seating plan: %w", err)
		}
		if res.RowsAffected() == 0 {
			return models.ErrNoSeatingPlan
		}
		return nil
	})
}

// SeatingPlanInput returns what a seating for all the remaining tours is
// made from: the participants to seat, the meetings and standings so far
// (SeatingHistory), and the first and last tour to seat: from the one after
// the last tour published, live or finished, to the tournament's last. It
// refuses a regular or finished event, numbers not drawn
// (models.ErrNotDrawn), participants that don't fill whole tables
// (models.ErrTablesIncomplete) and no tours left (models.ErrAllToursPlayed).
func (r *Repo) SeatingPlanInput(ctx context.Context, eventID int) (models.SeatingInput, int, int, error) {
	var (
		input       models.SeatingInput
		first, last int
	)
	err := r.InTx(ctx, func(tx *Repo) error {
		isTournament, _, err := tx.tournamentState(ctx, eventID, false)
		switch {
		case err != nil:
			return err
		case !isTournament:
			return models.ErrNotTournament
		}

		participants, err := tx.activeParticipants(ctx, eventID)
		if err != nil {
			return err
		}
		if missing := missingPlayers(len(participants)); missing > 0 {
			return models.WithData(models.ErrTablesIncomplete, models.ErrorData{"missing": missing})
		}

		rows, err := tx.db.Query(ctx, `
			SELECT
				COALESCE((SELECT max(t.number) FROM tours t
					WHERE t.event_id = $1 AND t.status IN ($2, $3, $4)), 0) AS seated,
				(SELECT (e.tournament->>'tours')::int FROM events e WHERE e.id = $1) AS tours`,
			eventID, models.TourPublished, models.TourLive, models.TourFinished)
		if err != nil {
			return fmt.Errorf("count tours: %w", err)
		}
		type counts struct {
			Seated int `db:"seated"`
			Tours  int `db:"tours"`
		}
		c, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[counts])
		if err != nil {
			return fmt.Errorf("count tours: %w", err)
		}
		first, last = c.Seated+1, c.Tours
		if first > last {
			return models.ErrAllToursPlayed
		}

		if input, err = tx.SeatingHistory(ctx, eventID, 0); err != nil {
			return err
		}
		input.Participants = participants
		return nil
	})
	return input, first, last, err
}
