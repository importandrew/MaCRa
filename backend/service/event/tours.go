package event

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/service/seating"
)

// GenerateTour drafts the tournament's next tour with tables made by the
// seating method (games with status created, seen by the staff only),
// replacing an earlier draft.
func (s *Service) GenerateTour(ctx context.Context, eventID int, how models.Seating, substitutes []int) (models.Tour, error) {
	makeTables, err := seatingFunc(how)
	if err != nil {
		return models.Tour{}, err
	}
	return s.repo.GenerateTour(ctx, eventID, how, substitutes, makeTables)
}

// PublishTour publishes the drafted seating to the players.
func (s *Service) PublishTour(ctx context.Context, eventID int) error {
	return s.repo.PublishTour(ctx, eventID)
}

// CancelTour cancels the published seating of the tour not started yet.
func (s *Service) CancelTour(ctx context.Context, eventID int) error {
	return s.repo.CancelTour(ctx, eventID)
}

// Tours returns the tournament's current tour with its tables and what the
// next tour would be, for its desk.
func (s *Service) Tours(ctx context.Context, eventID int) (models.TourState, error) {
	state, err := s.repo.TourState(ctx, eventID)
	if err != nil {
		return state, err
	}
	// The tables from a file, described; the client gets only their count.
	file := models.PlanTours{}
	for i := range state.Plan {
		if state.Plan[i].Tables != nil {
			file[state.Plan[i].Number] = state.Plan[i].Tables
		}
	}
	if len(file) > 0 {
		info := seating.FileInfo(file)
		state.File = &info
	}
	if state.Current == nil || len(state.Tables) == 0 {
		return state, nil
	}

	// The seating's quality, against the tours before it.
	input, err := s.repo.SeatingHistory(ctx, eventID, state.Current.Number)
	if err != nil {
		return state, err
	}
	tables := make([][4]int, 0, len(state.Tables))
	for _, table := range state.Tables {
		var seats [4]int
		for i, seat := range table.Seats {
			if i < len(seats) {
				seats[i] = seat.UserID
			}
			input.Participants = append(input.Participants, models.Participant{UserID: seat.UserID, Number: seat.Number})
		}
		tables = append(tables, seats)
	}
	quality, summary := seating.Quality(tables, input)
	for i := range state.Tables {
		state.Tables[i].Repeats = quality[i].Repeats
		state.Tables[i].RepeatPairs = quality[i].Pairs
		state.Tables[i].Places = quality[i].Places
	}
	state.Repeats, state.MaxMet, state.PlayersWithRepeats = summary.Repeats, summary.MaxMet, summary.PlayersWithRepeats
	return state, nil
}

// seatingSearchTime is how long a seating may search for fewer repeat
// meetings; planSearchTime the same for a random tour planned along with
// the random tours after it.
const (
	seatingSearchTime = 2 * time.Second
	planSearchTime    = 5 * time.Second
)

// checkSeating refuses an unknown method, a step with a method other than
// interval, and avoiding repeats with one other than Swiss or interval
// (models.ErrInvalidSeating).
func checkSeating(how models.Seating) error {
	switch {
	case !slices.Contains(models.ValidSeatingMethods, how.Method):
		return fmt.Errorf("%w: %q", models.ErrInvalidSeating, how.Method)
	case how.Step != nil && how.Method != models.SeatingInterval:
		return fmt.Errorf("%w: a step is for the interval seating only", models.ErrInvalidSeating)
	case how.AvoidRepeats && how.Method != models.SeatingSwiss && how.Method != models.SeatingInterval:
		return fmt.Errorf("%w: avoiding repeats is a choice of Swiss and interval only", models.ErrInvalidSeating)
	}
	return nil
}

// seatingFunc returns the table maker of a seating, or
// models.ErrInvalidSeating (checkSeating, also for an interval step below 1
// or above the number of tables).
func seatingFunc(how models.Seating) (models.TablesFunc, error) {
	if err := checkSeating(how); err != nil {
		return nil, err
	}
	switch how.Method {
	case models.SeatingRandom:
		return func(input models.SeatingInput) ([][4]int, error) {
			search := seatingSearchTime
			if input.FutureRandom > 0 {
				search = planSearchTime
			}
			return seating.Random(input, time.Now().Add(search))
		}, nil
	case models.SeatingSwiss:
		return func(input models.SeatingInput) ([][4]int, error) {
			return seating.Swiss(input, how.AvoidRepeats, time.Now().Add(seatingSearchTime))
		}, nil
	case models.SeatingInterval:
		return func(input models.SeatingInput) ([][4]int, error) {
			step := 0 // the number of tables
			if how.Step != nil {
				step = *how.Step
				if tables := len(input.Participants) / models.SeatsPerTable; step < 1 || step > tables {
					return nil, fmt.Errorf("%w: step %d, should be 1 to %d", models.ErrInvalidSeating, step, tables)
				}
			}
			return seating.Interval(input, step, how.AvoidRepeats, time.Now().Add(seatingSearchTime))
		}, nil
	default: // models.SeatingFixed
		return seating.Fixed, nil
	}
}

// SetTourPlan sets the seating a tour is to have. A step is for the
// interval method only, avoiding repeats for Swiss and interval (a random
// seating always avoids them, a fixed one never).
func (s *Service) SetTourPlan(ctx context.Context, eventID int, plan models.TourPlan) error {
	how := models.Seating{Method: plan.Seating, Step: plan.Step, AvoidRepeats: plan.AvoidRepeats}
	if err := checkSeating(how); err != nil {
		return err
	}
	if plan.Step != nil && *plan.Step < 1 {
		return fmt.Errorf("%w: step %d", models.ErrInvalidSeating, *plan.Step)
	}
	return s.repo.SetTourPlan(ctx, eventID, plan)
}

// ImportSeatingFile stores tables for several tours read from a CSV file
// (seating.ParseCSV), replacing the ones imported before; the tours in it
// are planned fixed.
func (s *Service) ImportSeatingFile(ctx context.Context, eventID int, file string) (models.SeatingFileInfo, error) {
	tours, err := seating.ParseCSV(file)
	if err != nil {
		return models.SeatingFileInfo{}, err
	}
	if err := s.repo.ImportSeatingFile(ctx, eventID, tours); err != nil {
		return models.SeatingFileInfo{}, err
	}
	return seating.FileInfo(tours), nil
}

// DeleteSeatingFile deletes the tables imported from a file.
func (s *Service) DeleteSeatingFile(ctx context.Context, eventID int) error {
	return s.repo.DeleteSeatingFile(ctx, eventID)
}

// StartTour starts the published tour: its games go live.
func (s *Service) StartTour(ctx context.Context, eventID int) error {
	return s.repo.StartTour(ctx, eventID)
}

// FinishTour finishes the live tour once all its games have ended.
func (s *Service) FinishTour(ctx context.Context, eventID int) error {
	return s.repo.FinishTour(ctx, eventID)
}

// AbortTour aborts the live tour: its games don't count.
func (s *Service) AbortTour(ctx context.Context, eventID int) error {
	return s.repo.AbortTour(ctx, eventID)
}

// AddTime adds minutes to a table of the live tour (a referee's action).
func (s *Service) AddTime(ctx context.Context, eventID, gameID, minutes int) error {
	if minutes < models.MinExtraMinutes || minutes > models.MaxExtraMinutes {
		return fmt.Errorf("%w: %d minutes, should be %d to %d",
			models.ErrInvalidExtraMinutes, minutes, models.MinExtraMinutes, models.MaxExtraMinutes)
	}
	return s.repo.AddTime(ctx, eventID, gameID, minutes)
}
