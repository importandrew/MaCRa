package event

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// AddMemberRoles adds roles to a member (or makes the user one). A finished
// event's members can't change (models.ErrEventAlreadyFinished).
func (s *Service) AddMemberRoles(ctx context.Context, eventID, targetUserID int, roles []models.EventRole) error {
	roles, err := memberRoles(roles)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return fmt.Errorf("%w: no roles to add", models.ErrInvalidRole)
	}

	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}
	if err := checkNotFinished(event); err != nil {
		return err
	}

	return s.repo.AddMemberRoles(ctx, eventID, targetUserID, roles)
}

// SetMemberRoles makes roles a member's roles (the role switcher in Manage →
// Members). The user must be a member already (models.ErrNotMember); add new
// members with AddMemberRoles. Only the owner may be left without roles
// (models.ErrInvalidRole for anyone else: remove them instead). The player
// role can't be taken from a member in a live game
// (models.ErrPlayerRoleInGame, with their name). A finished event's members
// can't change (models.ErrEventAlreadyFinished).
func (s *Service) SetMemberRoles(ctx context.Context, eventID, targetUserID int, roles []models.EventRole) error {
	roles, err := memberRoles(roles)
	if err != nil {
		return err
	}

	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}
	if err := checkNotFinished(event); err != nil {
		return err
	}

	isOwner := event.OwnerUserID == targetUserID
	current, err := s.repo.GetMemberRoles(ctx, eventID, targetUserID)
	if err != nil {
		return err
	}
	if len(current) == 0 && !isOwner {
		return fmt.Errorf("user %d: %w", targetUserID, models.ErrNotMember)
	}
	if len(roles) == 0 && !isOwner {
		return fmt.Errorf("%w: a member needs a role", models.ErrInvalidRole)
	}

	if slices.Contains(current, models.EventPlayer) && !slices.Contains(roles, models.EventPlayer) {
		if err := s.checkNotInGame(ctx, eventID, targetUserID, models.ErrPlayerRoleInGame); err != nil {
			return err
		}
	}

	return s.repo.SetMemberRoles(ctx, eventID, targetUserID, roles)
}

// checkNotInGame returns errInGame, with the user's name, if the user is in
// a live game of the event.
func (s *Service) checkNotInGame(ctx context.Context, eventID, userID int, errInGame error) error {
	_, err := s.gameService.CurrentGameID(ctx, userID, eventID)
	switch {
	case errors.Is(err, models.ErrGameNotFound):
		return nil
	case err != nil:
		return fmt.Errorf("user %d current game: %w", userID, err)
	default:
		return s.withPlayerName(ctx, userID, errInGame)
	}
}
