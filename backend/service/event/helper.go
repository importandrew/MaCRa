package event

import (
	"context"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) checkPlayerRole(ctx context.Context, eventID, playerID int) error {
	roles, err := s.repo.GetMemberRoles(ctx, eventID, playerID)
	if err != nil {
		return err
	}

	if !slices.Contains(roles, models.EventPlayer) {
		return models.ErrMissingPlayerRole
	}

	return nil
}
