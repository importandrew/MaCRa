package event

import (
	"context"
	"fmt"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/service/seating"
)

// GenerateTour drafts the tournament's next tour with tables made by the
// seating method (games with status created, seen by the staff only),
// replacing an earlier draft.
func (s *Service) GenerateTour(ctx context.Context, eventID int, method models.SeatingMethod) (models.Tour, error) {
	makeTables, err := seatingFunc(method)
	if err != nil {
		return models.Tour{}, err
	}
	return s.repo.GenerateTour(ctx, eventID, method, makeTables)
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
	return s.repo.TourState(ctx, eventID)
}

// seatingSearchTime is how long a seating may search for fewer repeat
// meetings.
const seatingSearchTime = 2 * time.Second

// seatingFunc returns the table maker of a seating method, or
// models.ErrInvalidSeating.
func seatingFunc(method models.SeatingMethod) (models.TablesFunc, error) {
	switch method {
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
			return seating.Interval(input, 0), nil
		}, nil
	default:
		return nil, fmt.Errorf("%w: %q", models.ErrInvalidSeating, method)
	}
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
