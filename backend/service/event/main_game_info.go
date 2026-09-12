package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) MainGameInfo(ctx context.Context, eventID, userID int) (models.MainGameInfo, error) {
	last, _ := s.gameService.LastGameResult(ctx, userID, eventID)
	now, _ := s.gameService.Now(ctx, userID, eventID)

	mainInfo := models.MainGameInfo{
		HasLast:    last.ID > 0,
		HasCurrent: now.ID > 0,
	}

	return mainInfo, nil
}
