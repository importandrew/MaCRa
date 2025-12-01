package game

import (
	"context"
	"fmt"
	"slices"

	"github.com/importandrew/MaCRa/engine/game"
)

func (s *Service) NewGame(ctx context.Context, playerIDs []int) (int, error) {
	if err := validatePlayers(playerIDs); err != nil {
		return 0, err
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

	g, err := game.NewGame(gameCfg)
	if err != nil {
		return 0, err
	}

	s.liveGames[g.ID()] = g

	return g.ID(), nil
}

func validatePlayers(playerIDs []int) error {
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

	return nil
}
