package event

import (
	"context"
	"fmt"
	"slices"

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
	if err := s.checkRolesFitEvent(ctx, eventID, roles); err != nil {
		return err
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
	if err := s.checkRolesFitEvent(ctx, eventID, roles); err != nil {
		return err
	}

	return s.repo.SetMemberRoles(ctx, eventID, targetUserID, roles)
}

// checkRolesFitEvent refuses the substitute role outside a tournament
// (models.ErrInvalidRole). That a member isn't both a player and a
// substitute is checked by the repo, with the roles they already hold.
func (s *Service) checkRolesFitEvent(ctx context.Context, eventID int, roles []models.EventRole) error {
	if !slices.Contains(roles, models.EventSubstitute) {
		return nil
	}
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}
	if event.Type != models.EventTournament {
		return fmt.Errorf("%w: substitutes are for tournaments only", models.ErrInvalidRole)
	}
	return nil
}
