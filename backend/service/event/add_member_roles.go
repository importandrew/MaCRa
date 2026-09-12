package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) AddMemberRoles(ctx context.Context, eventID, targetUserID int, roles []models.EventRole) error {
	return s.repo.AddMemberRoles(ctx, eventID, targetUserID, roles)
}
