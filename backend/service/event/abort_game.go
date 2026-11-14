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

func (s *Service) UndoLastRound(ctx context.Context, gameID, logSize int) error {
	return s.gameService.UndoLastRound(ctx, gameID, logSize)
}
