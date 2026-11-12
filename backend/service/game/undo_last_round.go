package game

import (
	"context"
	"time"

	"github.com/importandrew/MaCRa/models"
)

// UndoLastRound deletes the last round log entry of the game; an ended game
// becomes live again (within models.UndoFinalRoundTimeout). logSize is the
// number of log entries the referee saw; the repo refuses the undo with
// models.ErrLogChanged if the log has changed since.
func (s *Service) UndoLastRound(ctx context.Context, gameID, logSize int) error {
	return s.gameRepo.UndoLastRound(ctx, gameID, logSize, undoFinalRoundCutoff(time.Now()))
}

// undoFinalRoundCutoff returns the time a game must have ended after, as of
// now, for its final round to be undone.
func undoFinalRoundCutoff(now time.Time) time.Time {
	return now.Add(-models.UndoFinalRoundTimeout)
}
