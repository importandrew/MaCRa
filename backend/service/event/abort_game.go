package event

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

// AbortGame aborts a live game: a referee at any time, a seated player as the
// event's settings allow. When they allow it only while no round has been
// entered, that condition is checked again in the same statement as the
// abort, so a round entered at the same moment wins (models.ErrGameHasRounds).
func (s *Service) AbortGame(ctx context.Context, eventID, callerID, gameID int) error {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}

	authCtx, err := s.buildAuthContext(ctx, event, callerID, gameID)
	if err != nil {
		return err
	}

	if !authCtx.hasCapability(models.CapAbortGame) {
		return fmt.Errorf("user %d: %w", callerID, models.ErrEventForbidden)
	}

	onlyWithoutRounds := !authCtx.hasRole(models.EventReferee) && authCtx.seatedAbortOnlyWithoutRounds()
	return s.gameService.Abort(ctx, gameID, onlyWithoutRounds)
}

// UndoLastRound undoes the last round log entry of the game: a referee at any
// time (also the final round, within models.UndoFinalRoundTimeout, which
// makes an ended game live again), a seated player while the game is live if
// the event's settings allow it. The game may end between the check and the
// undo, so a seated player's undo refuses an ended game in the same
// transaction (models.ErrGameNotLive) instead of reviving it.
func (s *Service) UndoLastRound(ctx context.Context, eventID, callerID, gameID, logSize int) error {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}

	authCtx, err := s.buildAuthContext(ctx, event, callerID, gameID)
	if err != nil {
		return err
	}

	if !authCtx.hasCapability(models.CapUndoRound) {
		return fmt.Errorf("user %d: %w", callerID, models.ErrEventForbidden)
	}

	asReferee := authCtx.hasRole(models.EventReferee)
	return s.gameService.UndoLastRound(ctx, gameID, logSize, callerID, asReferee)
}
