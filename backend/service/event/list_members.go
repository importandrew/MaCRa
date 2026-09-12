package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) ListMembers(ctx context.Context, eventID int) ([]models.EventMember, error) {
	return s.repo.ListMembers(ctx, eventID)
}
