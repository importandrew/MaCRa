package event

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"slices"

	"github.com/importandrew/MaCRa/models"
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

// seatingFunc returns the table maker of a seating method, or
// models.ErrInvalidSeating.
func seatingFunc(method models.SeatingMethod) (models.TablesFunc, error) {
	switch method {
	case models.SeatingRandom:
		return randomTables, nil
	default:
		return nil, fmt.Errorf("%w: %q", models.ErrInvalidSeating, method)
	}
}

// randomTables shuffles the participants into tables of four. At a table the
// winds go by participant number, east the lowest (area 8: table and wind by
// participant number). The repo has checked that they fill whole tables.
func randomTables(participants []models.Participant) ([][4]int, error) {
	shuffled := slices.Clone(participants)
	if err := shuffle(shuffled); err != nil {
		return nil, err
	}

	tables := make([][4]int, 0, len(shuffled)/models.SeatsPerTable)
	for start := 0; start+models.SeatsPerTable <= len(shuffled); start += models.SeatsPerTable {
		table := slices.Clone(shuffled[start : start+models.SeatsPerTable])
		slices.SortFunc(table, func(a, b models.Participant) int { return a.Number - b.Number })
		tables = append(tables, [4]int{table[0].UserID, table[1].UserID, table[2].UserID, table[3].UserID})
	}
	return tables, nil
}

// shuffle puts the participants in a random order (Fisher–Yates) with
// crypto/rand: nobody should be able to predict the seating.
func shuffle(participants []models.Participant) error {
	for i := len(participants) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return fmt.Errorf("shuffle: %w", err)
		}
		participants[i], participants[j.Int64()] = participants[j.Int64()], participants[i]
	}
	return nil
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
