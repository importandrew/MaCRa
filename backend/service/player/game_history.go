package player

import (
	"context"
	"fmt"
)

func (s *Service) GameHistory(ctx context.Context, playerID int) ([]int, error) {
	_, err := s.repo.Get(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("player get: %w", err)
	}

	gameIDs, err := s.repo.GameHistory(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("player history: %w", err)
	}

	return gameIDs, nil
}
