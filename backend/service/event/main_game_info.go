package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) MainGameInfo(ctx context.Context, eventID, userID int) (models.MainGameInfo, error) {
	_, err := s.gameService.LastGameID(ctx, userID, eventID)
	hasLast, err := gameFound(err)
	if err != nil {
		return models.MainGameInfo{}, fmt.Errorf("last game: %w", err)
	}

	// A published tour's game counts: the player sees their table.
	_, err = s.gameService.SeatedGameID(ctx, userID, eventID)
	hasCurrent, err := gameFound(err)
	if err != nil {
		return models.MainGameInfo{}, fmt.Errorf("current game: %w", err)
	}

	return models.MainGameInfo{HasLast: hasLast, HasCurrent: hasCurrent}, nil
}

// gameFound turns the error of a game lookup into whether the game exists:
// models.ErrGameNotFound means it doesn't, any other error is passed on.
func gameFound(err error) (bool, error) {
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, models.ErrGameNotFound):
		return false, nil
	default:
		return false, err
	}
}
