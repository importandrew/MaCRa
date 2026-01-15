package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Status(ctx context.Context, gameID int) (models.GameStatus, error) {
	status, err := s.gameRepo.Status(ctx, gameID)
	if err != nil {
		return "", err
	}

	return status, nil
}
