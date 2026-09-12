package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) RemoveMembers(ctx context.Context, eventID int, targetUserIDs []int) error {
	for _, userID := range targetUserIDs {
		_, err := s.gameService.CurrentGameID(ctx, userID, eventID)
		if err != nil {
			if !errors.Is(err, models.ErrGameNotFound) {
				return fmt.Errorf("user %d current game: %w", userID, err)
			}
			continue
		}
		return fmt.Errorf("user %d: %w", userID, models.ErrMemberInCurrentGame)
	}

	return s.repo.RemoveMembers(ctx, eventID, targetUserIDs)
}

func (s *Service) GetMe(ctx context.Context, eventID, callerID int) (models.EventMe, error) {
	return s.MyCapabilities(ctx, eventID, callerID, 0)
}
