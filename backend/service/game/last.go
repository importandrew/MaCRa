package game

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) LastGameResult(ctx context.Context, playerID, eventID int) (models.GameResult, error) {
	gameID, err := s.gameRepo.LastGame(ctx, playerID, eventID)
	if err != nil {
		return models.GameResult{}, fmt.Errorf("find last game: %w", err)
	}

	res, err := s.GameResult(ctx, gameID)
	if err != nil {
		return models.GameResult{}, fmt.Errorf("game result: %w", err)
	}

	return res, nil
}

// LastGameID returns the ID of the player's last ended game in the event, or
// models.ErrGameNotFound if there is none.
func (s *Service) LastGameID(ctx context.Context, playerID, eventID int) (int, error) {
	return s.gameRepo.LastGame(ctx, playerID, eventID)
}
