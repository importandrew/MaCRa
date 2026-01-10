package game

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/importandrew/MaCRa/engine/game"
	"github.com/importandrew/MaCRa/models"
)

func (s *Service) NewGame(ctx context.Context, playerIDs []int) (int, error) {
	if err := s.validatePlayers(playerIDs); err != nil {
		return 0, fmt.Errorf("%w: %w", models.ErrValidatePlayers, err)
	}

	gameID, err := s.gameRepo.NewGame(ctx, [4]int(playerIDs))
	if err != nil {
		return 0, err
	}

	gameCfg := game.Config{
		ID:      gameID,
		Players: game.Seating(playerIDs),
		Rules:   game.Rules{},
	}

	g := game.NewGame(gameCfg)

	s.liveGames[gameID] = g

	return gameID, nil
}

func (s *Service) validatePlayers(playerIDs []int) error {
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
	for _, id := range playerIDs {
		g, err := s.CurrentGame(id)
		if err != nil {
			if !errors.Is(err, models.ErrGameNotFound) {
				// unexpected error
				return fmt.Errorf("current game: %w", err)
			}
			// ok, game not found
			continue
		}

		// game for player found
		return fmt.Errorf("player (ID %v) already in game (ID %v)", id, g.ID())
	}

	return nil
}
