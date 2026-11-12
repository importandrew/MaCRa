package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// AddMemberRoles adds roles to a member (or makes the user one). A finished
// event's members can't change (models.ErrEventAlreadyFinished).
func (s *Service) AddMemberRoles(ctx context.Context, eventID, targetUserID int, roles []models.EventRole) error {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}
	if err := checkNotFinished(event); err != nil {
		return err
	}

	return s.repo.AddMemberRoles(ctx, eventID, targetUserID, roles)
}
