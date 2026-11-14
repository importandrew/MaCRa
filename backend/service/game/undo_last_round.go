package game

import (
	"context"
	"time"

	"github.com/importandrew/MaCRa/models"
)

// UndoLastRound deletes the last round log entry of the game. With
// reviveEnded an ended game becomes live again (within
// models.UndoFinalRoundTimeout); without it only a live game is undone
// (models.ErrGameNotLive). logSize is the number of log entries the caller
// saw; the repo refuses the undo with models.ErrLogChanged if the log has
// changed since.
func (s *Service) UndoLastRound(ctx context.Context, gameID, logSize int, reviveEnded bool) error {
	var endedAfter *time.Time
	if reviveEnded {
		endedAfter = new(UndoFinalRoundCutoff(time.Now()))
	}
	return s.gameRepo.UndoLastRound(ctx, gameID, logSize, endedAfter)
}

// UndoFinalRoundCutoff returns the time a game must have ended after, as of
// now, for its final round to be undone.
func UndoFinalRoundCutoff(now time.Time) time.Time {
	return now.Add(-models.UndoFinalRoundTimeout)
}
