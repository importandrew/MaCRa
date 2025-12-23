package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Abort(ctx context.Context, gameID int) error {
	_, err := s.GetLiveGame(gameID)
	if err != nil {
		return err
	}

	err = s.gameRepo.SetStatus(ctx, gameID, models.GameAborted)
	if err != nil {
		return err
	}

	delete(s.liveGames, gameID)

	return nil
}
