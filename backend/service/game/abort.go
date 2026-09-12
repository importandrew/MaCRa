package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Abort(ctx context.Context, gameID int) error {
	status, err := s.gameRepo.Status(ctx, gameID)
	if err != nil {
		return err
	}

	if status != models.GameLive {
		return models.ErrGameNotLive
	}

	err = s.gameRepo.SetStatus(ctx, gameID, models.GameAborted)
	if err != nil {
		return err
	}

	return nil
}
