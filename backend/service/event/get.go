package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) Get(ctx context.Context, eventID int) (models.Event, error) {
	return s.repo.Get(ctx, eventID)
}

func (s *Service) ListVisible(ctx context.Context, userID int) ([]models.EventListItem, error) {
	events, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	filtered := make([]models.EventListItem, 0, len(events))
	for _, event := range events {
		authCtx, err := s.buildAuthContext(ctx, event, userID, 0)
		if err != nil {
			return nil, err
		}

		if !authCtx.canAccessEvent() {
			continue
		}

		item := models.EventListItem{Event: event}
		if userID != 0 {
			item.IsParticipating = authCtx.isParticipating()
			item.IsOwner = authCtx.IsOwner
			item.Roles = authCtx.Roles
		}

		filtered = append(filtered, item)
	}

	return filtered, nil
}
