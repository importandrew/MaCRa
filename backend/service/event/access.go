package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) CanAccess(ctx context.Context, eventID, callerID int) error {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}

	authCtx, err := s.buildAuthContext(ctx, event, callerID, 0)
	if err != nil {
		return err
	}

	if !authCtx.canAccessEvent() {
		return models.ErrEventForbidden
	}

	return nil
}
