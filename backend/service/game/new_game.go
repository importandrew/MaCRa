package game

import (
	"context"
	"fmt"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// NewGame creates a live game with the seating (east, south, west, north).
// Whether the players belong to the event is checked in the event service;
// that none of them is in another live game, by the repo under its locks
// (models.ErrPlayerInAnotherGame).
func (s *Service) NewGame(ctx context.Context, playerIDs []int, eventID int) (int, error) {
	if err := validatePlayers(playerIDs); err != nil {
		return 0, fmt.Errorf("%w: %w", models.ErrValidatePlayers, err)
	}

	return s.gameRepo.NewGame(ctx, [4]int(playerIDs), eventID)
}

// validatePlayers checks the seating: four distinct players.
func validatePlayers(playerIDs []int) error {
	if len(playerIDs) != 4 {
		return fmt.Errorf("expected four player IDs, got %v", len(playerIDs))
	}

	if slices.Contains(playerIDs, 0) {
		return fmt.Errorf("missing player at index %d", slices.Index(playerIDs, 0))
	}

	sorted := slices.Clone(playerIDs)
	slices.Sort(sorted)
	if len(slices.Compact(sorted)) != 4 {
		return fmt.Errorf("expected four unique players, got ids %v", playerIDs)
	}

	return nil
}
