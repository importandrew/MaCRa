package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// RemoveMembers removes members from the event. The repo refuses, under the
// same locks as starting a game, the owner (models.ErrOwnerNotRemovable), a
// member in a live game (models.ErrMemberInCurrentGame) and a finished event
// (models.ErrEventAlreadyFinished).
func (s *Service) RemoveMembers(ctx context.Context, eventID int, targetUserIDs []int) error {
	return s.repo.RemoveMembers(ctx, eventID, targetUserIDs)
}

func (s *Service) GetMe(ctx context.Context, eventID, callerID int) (models.EventMe, error) {
	return s.MyCapabilities(ctx, eventID, callerID)
}
