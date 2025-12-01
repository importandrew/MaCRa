package game

import (
	"context"
	"errors"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Abort(ctx context.Context, gameID int) error {
	_, err := s.GetLiveGame(gameID)
	if err != nil {
		if errors.Is(err, errGameNotFound) {
			return errors.New("game is not live")
		}
		return err
	}

	err = s.gameRepo.SetStatus(ctx, gameID, models.GameAborted)
	if err != nil {
		return err
	}

	delete(s.liveGames, gameID)

	return nil
}
