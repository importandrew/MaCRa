package game

import "context"

func (s *Service) EventID(ctx context.Context, gameID int) (int, error) {
	return s.gameRepo.EventID(ctx, gameID)
}
