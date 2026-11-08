package game

import (
	"context"
)

// UndoLastRound deletes the last round log entry of the game; an ended game
// becomes live again (within models.UndoFinalRoundTimeout). logSize is the
// number of log entries the referee saw; the repo refuses the undo with
// models.ErrLogChanged if the log has changed since.
func (s *Service) UndoLastRound(ctx context.Context, gameID, logSize int) error {
	return s.gameRepo.UndoLastRound(ctx, gameID, logSize)
}
