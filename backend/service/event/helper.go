package event

import (
	"context"
	"fmt"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

// checkPlayerRole returns models.ErrMissingPlayerRole, with the player's name
// for the client, unless the user has the player role in the event.
func (s *Service) checkPlayerRole(ctx context.Context, eventID, playerID int) error {
	roles, err := s.repo.GetMemberRoles(ctx, eventID, playerID)
	if err != nil {
		return err
	}

	if !slices.Contains(roles, models.EventPlayer) {
		return s.withPlayerName(ctx, playerID, models.ErrMissingPlayerRole)
	}

	return nil
}

// withPlayerName returns err with the player's name for the client and the ID
// for the log.
func (s *Service) withPlayerName(ctx context.Context, playerID int, err error) error {
	player, getErr := s.playerRepo.Get(ctx, playerID)
	if getErr != nil {
		return fmt.Errorf("player %d: %w (name: %w)", playerID, err, getErr)
	}

	return fmt.Errorf("player %d: %w", playerID,
		models.WithData(err, models.ErrorData{"name": player.Name}))
}
