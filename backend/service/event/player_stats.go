package event

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// PlayerStats returns the player's statistics in the event, from its ended
// games.
func (s *Service) PlayerStats(ctx context.Context, eventID, playerID int) (models.PlayerStats, error) {
	return s.repo.PlayerStats(ctx, eventID, playerID)
}
