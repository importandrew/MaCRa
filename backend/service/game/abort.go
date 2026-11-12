package game

import (
	"context"
)

// Abort makes a live game aborted: it won't count. With onlyWithoutRounds
// (a seated player, not a referee) only while no round has been entered,
// else models.ErrGameHasRounds; checked by the repo in the same statement.
func (s *Service) Abort(ctx context.Context, gameID int, onlyWithoutRounds bool) error {
	return s.gameRepo.Abort(ctx, gameID, onlyWithoutRounds)
}
