package game

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) UndoLastRound(ctx context.Context, gameID int) error {
	status, err := s.gameRepo.Status(ctx, gameID)
	if err != nil {
		return err
	}

	switch status {
	case models.GameLive:
		if err := s.gameRepo.DeleteLastRound(ctx, gameID); err != nil {
			return err
		}
	case models.GameEnded:
		if err := s.gameRepo.Unfinalize(ctx, gameID); err != nil {
			return err
		}
		if err := s.gameRepo.DeleteLastRound(ctx, gameID); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected game status: %v", status)
	}

	return nil
}
