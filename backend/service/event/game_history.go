package event

import (
	"context"
)

func (s *Service) GameHistory(ctx context.Context, eventID int) ([]int, error) {
	return s.repo.GameHistory(ctx, eventID)
}

func (s *Service) PlayerGameHistory(ctx context.Context, eventID, playerID int) ([]int, error) {
	return s.repo.PlayerGameHistory(ctx, eventID, playerID)
}
