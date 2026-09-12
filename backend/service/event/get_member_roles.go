package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) GetMemberRoles(ctx context.Context, eventID, callerID, targetUserID int) ([]models.EventRole, error) {
	return s.repo.GetMemberRoles(ctx, eventID, targetUserID)
}
