package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// ListMembers returns the event's members with their roles, for anyone who
// can open the event (no emails: they are never shown to other users).
func (s *Service) ListMembers(ctx context.Context, eventID int) ([]models.EventMember, error) {
	return s.repo.ListMembers(ctx, eventID)
}
