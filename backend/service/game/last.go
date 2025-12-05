package game

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) LastGameResult(ctx context.Context, playerID int) (models.GameResult, error) {
	gameID, err := s.gameRepo.LastGame(ctx, playerID)
	if err != nil {
		return models.GameResult{}, fmt.Errorf("find last game: %w", err)
	}

	res, err := s.GameResult(ctx, gameID)
	if err != nil {
		return models.GameResult{}, fmt.Errorf("game result: %w", err)
	}

	return res, nil
}
