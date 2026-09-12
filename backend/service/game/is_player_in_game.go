package game

import "context"

func (s *Service) IsPlayerInGame(ctx context.Context, gameID, playerID int) (bool, error) {
	return s.gameRepo.IsPlayerInGame(ctx, gameID, playerID)
}
