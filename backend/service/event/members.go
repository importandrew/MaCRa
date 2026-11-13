package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

// RemoveMembers removes members from the event: not while one of them is in
// a live game (models.ErrMemberInCurrentGame), and not in a finished event
// (models.ErrEventAlreadyFinished), whose members can't change.
func (s *Service) RemoveMembers(ctx context.Context, eventID int, targetUserIDs []int) error {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}
	if err := checkNotFinished(event); err != nil {
		return err
	}

	for _, userID := range targetUserIDs {
		_, err := s.gameService.CurrentGameID(ctx, userID, eventID)
		if err != nil {
			if !errors.Is(err, models.ErrGameNotFound) {
				return fmt.Errorf("user %d current game: %w", userID, err)
			}
			continue
		}
		return s.withPlayerName(ctx, userID, models.ErrMemberInCurrentGame)
	}

	return s.repo.RemoveMembers(ctx, eventID, targetUserIDs)
}

func (s *Service) GetMe(ctx context.Context, eventID, callerID int) (models.EventMe, error) {
	return s.MyCapabilities(ctx, eventID, callerID)
}
