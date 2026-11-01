package event

import (
	"context"
	"fmt"

	"github.com/importandrew/MaCRa/models"
)

func (s *Service) HasCapability(ctx context.Context, eventID, callerID, gameID int, capability models.Capability) error {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return err
	}

	authCtx, err := s.buildAuthContext(ctx, event, callerID, gameID)
	if err != nil {
		return err
	}

	if !authCtx.hasCapability(capability) {
		return fmt.Errorf("user %d: %w", callerID, models.ErrEventForbidden)
	}

	return nil
}

func (s *Service) MyCapabilities(ctx context.Context, eventID, callerID, gameID int) (models.EventMe, error) {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return models.EventMe{}, err
	}

	authCtx, err := s.buildAuthContext(ctx, event, callerID, gameID)
	if err != nil {
		return models.EventMe{}, err
	}

	allCaps := []models.Capability{
		models.CapManageEvent,
		models.CapManageMembers,
		models.CapCreateGame,
		models.CapAbortGame,
		models.CapFinishGame,
		models.CapUndoRound,
		models.CapWatchLive,
		models.CapWatchResults,
		models.CapViewRatings,
	}

	caps := make([]models.Capability, 0, len(allCaps))
	for _, capability := range allCaps {
		if authCtx.hasCapability(capability) {
			caps = append(caps, capability)
		}
	}

	me := models.EventMe{
		IsOwner:      authCtx.IsOwner,
		Roles:        authCtx.Roles,
		Capabilities: caps,
	}

	return me, nil
}
