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

// MyCapabilities returns the caller's event-level capabilities
// (models.EventCapabilities) and roles.
func (s *Service) MyCapabilities(ctx context.Context, eventID, callerID int) (models.EventMe, error) {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return models.EventMe{}, err
	}

	authCtx, err := s.buildAuthContext(ctx, event, callerID, 0)
	if err != nil {
		return models.EventMe{}, err
	}

	me := models.EventMe{
		IsOwner:      authCtx.IsOwner,
		Roles:        authCtx.Roles,
		Capabilities: authCtx.allowed(models.EventCapabilities),
	}
	if me.JoinRequest, err = s.joinRequestState(ctx, authCtx); err != nil {
		return models.EventMe{}, err
	}
	if authCtx.hasCapability(models.CapManageMembers) && event.FinishedAt == nil {
		if me.PendingRequests, err = s.repo.CountJoinRequests(ctx, eventID); err != nil {
			return models.EventMe{}, err
		}
	}
	return me, nil
}

// MyGameCapabilities returns the caller's capabilities for one game
// (models.GameCapabilities).
func (s *Service) MyGameCapabilities(ctx context.Context, eventID, callerID, gameID int) (models.GameMe, error) {
	event, err := s.repo.Get(ctx, eventID)
	if err != nil {
		return models.GameMe{}, err
	}

	authCtx, err := s.buildAuthContext(ctx, event, callerID, gameID)
	if err != nil {
		return models.GameMe{}, err
	}

	return models.GameMe{Capabilities: authCtx.allowed(models.GameCapabilities)}, nil
}

// allowed returns the capabilities of caps the caller has.
func (authCtx EventAuthContext) allowed(caps []models.Capability) []models.Capability {
	out := make([]models.Capability, 0, len(caps))
	for _, capability := range caps {
		if authCtx.hasCapability(capability) {
			out = append(out, capability)
		}
	}
	return out
}
