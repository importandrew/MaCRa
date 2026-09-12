package game

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// NewGame checks if players do not currently participate in game in this event.
// Whether players belong to the event is checked in event service.
func (s *Service) NewGame(ctx context.Context, playerIDs []int, eventID int) (int, error) {
	if err := s.validatePlayers(ctx, playerIDs, eventID); err != nil {
		return 0, fmt.Errorf("%w: %w", models.ErrValidatePlayers, err)
	}

	gameID, err := s.gameRepo.NewGame(ctx, [4]int(playerIDs), eventID)
	if err != nil {
		return 0, err
	}

	// gameCfg := game.Config{
	// 	ID:      gameID,
	// 	Players: game.Seating(playerIDs),
	// 	Rules:   game.Rules{},
	// }

	// g := game.NewGame(gameCfg)

	// _, ok := s.liveGames[eventID]
	// if !ok {
	// 	s.liveGames[eventID] = make(map[int]GameEngine)
	// }
	// s.liveGames[eventID][gameID] = g

	return gameID, nil
}

// currentGameID should be nil unless it's game reload (otherwise interferes with itself).
func (s *Service) validatePlayers(ctx context.Context, playerIDs []int, eventID int) error {
	// len = 4
	if len(playerIDs) != 4 {
		return fmt.Errorf("expected four player IDs, got %v", len(playerIDs))
	}

	// not empty
	if slices.Contains(playerIDs, 0) {
		return fmt.Errorf("missing player at index %d", slices.Index(playerIDs, 0))
	}

	// uniqueness
	sorted := make([]int, len(playerIDs))
	copy(sorted, playerIDs)

	slices.Sort(sorted)

	if len(slices.Compact(sorted)) != 4 {
		return fmt.Errorf("expected four unique players, got ids %v", playerIDs)
	}

	// disallow same player in several games at the same time
	for _, playerID := range playerIDs {
		gameID, err := s.gameRepo.CurrentGame(ctx, playerID, eventID)
		if err != nil {
			if !errors.Is(err, models.ErrGameNotFound) {
				// unexpected error
				return fmt.Errorf("current game: %w", err)
			}
			// ok, game not found
			continue
		}

		return fmt.Errorf("player (ID %v) already in game (ID %v)", playerID, gameID)
	}

	return nil
}
