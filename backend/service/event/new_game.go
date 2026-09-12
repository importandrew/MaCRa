package event

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

// NewGame checks if players belong to the event and calls game service to create new game.
// Other validations are checked inside game service.
func (s *Service) NewGame(ctx context.Context, eventID int, playerIDs []int) (int, error) {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return 0, err
	}

	if event.FinishedAt != nil {
		return 0, models.ErrEventAlreadyFinished
	}

	for _, playerID := range playerIDs {
		if err := s.checkPlayerRole(ctx, eventID, playerID); err != nil {
			return 0, fmt.Errorf("player ID %d: %w", playerID, err)
		}
	}

	return s.gameService.NewGame(ctx, playerIDs, eventID)
}
