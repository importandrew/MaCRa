package event

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// PublishTour publishes the tournament's next tour with tables made by the
// seating method, creating its games (status created).
func (s *Service) PublishTour(ctx context.Context, eventID int, method models.SeatingMethod) (models.Tour, error) {
	makeTables, err := seatingFunc(method)
	if err != nil {
		return models.Tour{}, err
	}
	return s.repo.PublishTour(ctx, eventID, method, makeTables)
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
