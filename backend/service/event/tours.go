package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/service/seating"
)

// GenerateTour drafts the tournament's next tour with tables made by the
// seating method (games with status created, seen by the staff only),
// replacing an earlier draft.
func (s *Service) GenerateTour(ctx context.Context, eventID int, how models.Seating) (models.Tour, error) {
	makeTables, err := seatingFunc(how)
	if err != nil {
		return models.Tour{}, err
	}
	return s.repo.GenerateTour(ctx, eventID, how, makeTables)
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
	plan, err := s.repo.SeatingPlan(ctx, eventID)
	switch {
	case errors.Is(err, models.ErrNoSeatingPlan):
	case err != nil:
		return state, err
	default:
		info := seating.PlanInfo(plan)
		state.Plan = &info
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
// meetings.
const seatingSearchTime = 2 * time.Second

// seatingFunc returns the table maker of a seating method, or
// models.ErrInvalidSeating (also for an interval step below 1 or above the
// number of tables, or a step with another method).
func seatingFunc(how models.Seating) (models.TablesFunc, error) {
	if how.Step != nil && how.Method != models.SeatingInterval {
		return nil, fmt.Errorf("%w: a step is for the interval seating only", models.ErrInvalidSeating)
	}
	switch how.Method {
	case models.SeatingRandom:
		return func(input models.SeatingInput) ([][4]int, error) {
			return seating.Random(input, time.Now().Add(seatingSearchTime))
		}, nil
	case models.SeatingSwiss:
		return func(input models.SeatingInput) ([][4]int, error) {
			return seating.Swiss(input, time.Now().Add(seatingSearchTime))
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
			return seating.Interval(input, step), nil
		}, nil
	case models.SeatingFixed:
		return seating.Fixed, nil
	default:
		return nil, fmt.Errorf("%w: %q", models.ErrInvalidSeating, how.Method)
	}
}

// planSearchTime is how long a seating for all tours may search for fewer
// repeat meetings.
const planSearchTime = 5 * time.Second

// GenerateSeatingPlan prepares the seating for all the remaining tours at
// once (from the one after the last published), with as few repeat meetings
// as the search finds, replacing an earlier plan.
func (s *Service) GenerateSeatingPlan(ctx context.Context, eventID int) (models.SeatingPlanInfo, error) {
	input, first, last, err := s.repo.SeatingPlanInput(ctx, eventID)
	if err != nil {
		return models.SeatingPlanInfo{}, err
	}
	tours, err := seating.Plan(input, last-first+1, time.Now().Add(planSearchTime))
	if err != nil {
		return models.SeatingPlanInfo{}, err
	}
	plan := models.PlanTours{}
	for i, tables := range tours {
		plan[first+i] = tables
	}
	return s.savePlan(ctx, eventID, models.PlanGenerated, plan)
}

// ImportSeatingPlan stores a seating for several tours read from a CSV file
// (seating.ParseCSV), replacing an earlier plan.
func (s *Service) ImportSeatingPlan(ctx context.Context, eventID int, file string) (models.SeatingPlanInfo, error) {
	plan, err := seating.ParseCSV(file)
	if err != nil {
		return models.SeatingPlanInfo{}, err
	}
	return s.savePlan(ctx, eventID, models.PlanImported, plan)
}

// savePlan stores a seating plan and describes it.
func (s *Service) savePlan(ctx context.Context, eventID int, source models.SeatingPlanSource, plan models.PlanTours) (models.SeatingPlanInfo, error) {
	if err := s.repo.SaveSeatingPlan(ctx, eventID, source, plan); err != nil {
		return models.SeatingPlanInfo{}, err
	}
	return seating.PlanInfo(models.SeatingPlan{Source: source, CreatedAt: time.Now().UTC(), Tours: plan}), nil
}

// DeleteSeatingPlan deletes the seating prepared for all tours.
func (s *Service) DeleteSeatingPlan(ctx context.Context, eventID int) error {
	return s.repo.DeleteSeatingPlan(ctx, eventID)
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
