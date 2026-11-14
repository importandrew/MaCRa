package event

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

// AddMemberRoles adds roles to a member (or makes the user one): at least one
// known role (models.ErrInvalidRole). The repo refuses a finished event
// (models.ErrEventAlreadyFinished).
func (s *Service) AddMemberRoles(ctx context.Context, eventID, targetUserID int, roles []models.EventRole) error {
	roles, err := memberRoles(roles)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return fmt.Errorf("%w: no roles to add", models.ErrInvalidRole)
	}

	return s.repo.AddMemberRoles(ctx, eventID, targetUserID, roles)
}

// SetMemberRoles makes roles a member's roles (the role switcher in Manage →
// Members): known roles, each once (models.ErrInvalidRole). The repo checks
// the rest under the same locks as starting a game: the user is a member,
// only the owner may be left without roles, the player role stays while the
// member is in a live game, the event isn't finished.
func (s *Service) SetMemberRoles(ctx context.Context, eventID, targetUserID int, roles []models.EventRole) error {
	roles, err := memberRoles(roles)
	if err != nil {
		return err
	}

	return s.repo.SetMemberRoles(ctx, eventID, targetUserID, roles)
}
