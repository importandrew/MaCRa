package game

import (
	"context"

	"github.com/importandrew/MaCRa/models"
)

// Summary returns the game's status, end time and log size, for the
// capability rules.
func (s *Service) Summary(ctx context.Context, gameID int) (models.GameSummary, error) {
	return s.gameRepo.Summary(ctx, gameID)
}
