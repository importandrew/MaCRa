package event

import (
	"context"
	"fmt"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// RemoveMembers removes members from the event: not the owner
// (models.ErrOwnerNotRemovable; they may drop their roles instead), not
// while one of them is in a live game (models.ErrMemberInCurrentGame), and
// not in a finished event (models.ErrEventAlreadyFinished), whose members
// can't change.
func (s *Service) RemoveMembers(ctx context.Context, eventID int, targetUserIDs []int) error {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}
	if err := checkNotFinished(event); err != nil {
		return err
	}
	if slices.Contains(targetUserIDs, event.OwnerUserID) {
		return fmt.Errorf("user %d: %w", event.OwnerUserID, models.ErrOwnerNotRemovable)
	}

	for _, userID := range targetUserIDs {
		if err := s.checkNotInGame(ctx, eventID, userID, models.ErrMemberInCurrentGame); err != nil {
			return err
		}
	}

	return s.repo.RemoveMembers(ctx, eventID, targetUserIDs)
}

func (s *Service) GetMe(ctx context.Context, eventID, callerID int) (models.EventMe, error) {
	return s.MyCapabilities(ctx, eventID, callerID)
}
