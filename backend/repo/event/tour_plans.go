package event

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/importandrew/MaCRa/models"
)

// dataTour is the key of a tour's number in an error's data.
const dataTour = "tour"

// tourPlans returns the tournament's plan rows by tour number (tours
// without a row are random) and its number of tours.
func (r *Repo) tourPlans(ctx context.Context, eventID int) (map[int]models.TourPlan, int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT number, seating, seating_step, avoid_repeats, tables
		FROM tour_plans WHERE event_id = $1`, eventID)
	if err != nil {
		return nil, 0, fmt.Errorf("tour plans: %w", err)
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.TourPlan])
	if err != nil {
		return nil, 0, fmt.Errorf("tour plans: %w", err)
	}
	plans := make(map[int]models.TourPlan, len(list))
	for _, plan := range list {
		plans[plan.Number] = plan
	}

	rows, err = r.db.Query(ctx, `
		SELECT COALESCE((tournament->>'tours')::int, 0) FROM events WHERE id = $1`, eventID)
	if err != nil {
		return nil, 0, fmt.Errorf("number of tours: %w", err)
	}
	tours, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return nil, 0, fmt.Errorf("number of tours: %w", err)
	}
	return plans, tours, nil
}

// planList returns the plan of each tour, 1 to the number of tours.
func (r *Repo) planList(ctx context.Context, eventID int) ([]models.TourPlan, error) {
	plans, tours, err := r.tourPlans(ctx, eventID)
	if err != nil {
		return nil, err
	}
	seated, err := r.seatedTours(ctx, eventID)
	if err != nil {
		return nil, err
	}
	list := make([]models.TourPlan, 0, tours)
	for n := 1; n <= tours; n++ {
		plan, ok := plans[n]
		if !ok {
			plan = models.TourPlan{Number: n, Seating: models.SeatingRandom}
		}
		plan.FileTables = len(plan.Tables)
		plan.Seated = n <= seated
		list = append(list, plan)
	}
	return list, nil
}

// planInput returns what the plan gives tour n's seating: its tables from a
// file (nil without them) and the number of random tours after it.
func (r *Repo) planInput(ctx context.Context, eventID, n int) ([][4]int, int, error) {
	plans, tours, err := r.tourPlans(ctx, eventID)
	if err != nil {
		return nil, 0, err
	}
	futureRandom := 0
	for after := n + 1; after <= tours; after++ {
		if plan, ok := plans[after]; !ok || plan.Seating == models.SeatingRandom {
			futureRandom++
		}
	}
	return plans[n].Tables, futureRandom, nil
}

// seatedTours returns the number of the last tour published, live or
// finished (0 before the first): the plan of the tours up to it no longer
// matters.
func (r *Repo) seatedTours(ctx context.Context, eventID int) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT COALESCE(max(number), 0) FROM tours
		WHERE event_id = $1 AND status IN ($2, $3, $4)`,
		eventID, models.TourPublished, models.TourLive, models.TourFinished)
	if err != nil {
		return 0, fmt.Errorf("seated tours: %w", err)
	}
	seated, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, fmt.Errorf("seated tours: %w", err)
	}
	return seated, nil
}

// lockTournament locks the event row FOR UPDATE and refuses a regular
// (models.ErrNotTournament) or finished event.
func (r *Repo) lockTournament(ctx context.Context, eventID int) error {
	isTournament, _, err := r.tournamentState(ctx, eventID, true)
	switch {
	case err != nil:
		return err
	case !isTournament:
		return models.ErrNotTournament
	}
	return nil
}

// SetTourPlan sets the seating a tour is to have, keeping its tables from a
// file. It refuses a tour out of the tournament's range
// (models.ErrInvalidSeating), one already seated (models.ErrTourAlreadySeated)
// and the fixed method for a tour without tables from a file
// (models.ErrFileSeatingMissing).
func (r *Repo) SetTourPlan(ctx context.Context, eventID int, plan models.TourPlan) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if err := tx.lockTournament(ctx, eventID); err != nil {
			return err
		}
		plans, tours, err := tx.tourPlans(ctx, eventID)
		if err != nil {
			return err
		}
		seated, err := tx.seatedTours(ctx, eventID)
		if err != nil {
			return err
		}
		switch {
		case plan.Number < 1 || plan.Number > tours:
			return fmt.Errorf("%w: tour %d of %d", models.ErrInvalidSeating, plan.Number, tours)
		case plan.Number <= seated:
			return models.WithData(models.ErrTourAlreadySeated, models.ErrorData{dataTour: plan.Number})
		case plan.Seating == models.SeatingFixed && plans[plan.Number].Tables == nil:
			return models.WithData(models.ErrFileSeatingMissing, models.ErrorData{dataTour: plan.Number})
		}

		_, err = tx.db.Exec(ctx, `
			INSERT INTO tour_plans (event_id, number, seating, seating_step, avoid_repeats)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (event_id, number) DO UPDATE SET seating = EXCLUDED.seating,
				seating_step = EXCLUDED.seating_step, avoid_repeats = EXCLUDED.avoid_repeats`,
			eventID, plan.Number, plan.Seating, plan.Step, plan.AvoidRepeats)
		if err != nil {
			return fmt.Errorf("set tour plan: %w", err)
		}
		return nil
	})
}

// ImportSeatingFile stores tables imported from a file, replacing the ones
// imported before. The tours in the file not seated yet are planned fixed;
// tours planned fixed that the file doesn't have become random. It refuses
// a tour beyond the tournament's tours (models.ErrSeatingFileTour).
func (r *Repo) ImportSeatingFile(ctx context.Context, eventID int, file models.PlanTours) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if err := tx.lockTournament(ctx, eventID); err != nil {
			return err
		}
		_, tours, err := tx.tourPlans(ctx, eventID)
		if err != nil {
			return err
		}
		for tour := range file {
			if tour > tours {
				return models.WithData(models.ErrSeatingFileTour, models.ErrorData{dataTour: tour, "tours": tours})
			}
		}
		seated, err := tx.seatedTours(ctx, eventID)
		if err != nil {
			return err
		}

		if err := tx.dropFileTables(ctx, eventID, seated); err != nil {
			return err
		}
		for tour, tables := range file {
			_, err := tx.db.Exec(ctx, `
				INSERT INTO tour_plans (event_id, number, seating, avoid_repeats, tables)
				VALUES ($1, $2, $3, FALSE, $4)
				ON CONFLICT (event_id, number) DO UPDATE SET tables = EXCLUDED.tables,
					seating = CASE WHEN tour_plans.number > $5 THEN EXCLUDED.seating ELSE tour_plans.seating END,
					seating_step = CASE WHEN tour_plans.number > $5 THEN NULL ELSE tour_plans.seating_step END,
					avoid_repeats = CASE WHEN tour_plans.number > $5 THEN FALSE ELSE tour_plans.avoid_repeats END`,
				eventID, tour, models.SeatingFixed, tables, seated)
			if err != nil {
				return fmt.Errorf("import tour %d: %w", tour, err)
			}
		}
		return nil
	})
}

// DeleteSeatingFile deletes the tables imported from a file; the tours not
// seated yet that were planned fixed become random.
func (r *Repo) DeleteSeatingFile(ctx context.Context, eventID int) error {
	return r.InTx(ctx, func(tx *Repo) error {
		if err := tx.lockTournament(ctx, eventID); err != nil {
			return err
		}
		seated, err := tx.seatedTours(ctx, eventID)
		if err != nil {
			return err
		}
		return tx.dropFileTables(ctx, eventID, seated)
	})
}

// dropFileTables forgets the tables from a file; the tours after seated that
// were planned fixed become random.
func (r *Repo) dropFileTables(ctx context.Context, eventID, seated int) error {
	if _, err := r.db.Exec(ctx, `
		UPDATE tour_plans SET seating = $3, seating_step = NULL, avoid_repeats = FALSE
		WHERE event_id = $1 AND number > $2 AND seating = $4`,
		eventID, seated, models.SeatingRandom, models.SeatingFixed); err != nil {
		return fmt.Errorf("unplan fixed tours: %w", err)
	}
	if _, err := r.db.Exec(ctx, `
		UPDATE tour_plans SET tables = NULL WHERE event_id = $1`, eventID); err != nil {
		return fmt.Errorf("drop file tables: %w", err)
	}
	return nil
}
