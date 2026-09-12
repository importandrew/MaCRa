package event

import (
	"context"
	"slices"

	"github.com/importandrew/MaCRa/models"
)

type EventAuthContext struct {
	Event          models.Event
	CallerID       int
	Roles          []models.EventRole
	IsOwner        bool
	IsMember       bool
	IsSeatedInGame bool
}

func (s *Service) buildAuthContext(ctx context.Context, event models.Event, callerID, gameID int) (EventAuthContext, error) {
	authCtx := EventAuthContext{
		Event:    event,
		CallerID: callerID,
		IsOwner:  event.OwnerUserID == callerID,
	}

	if callerID != 0 {
		roles, err := s.repo.GetMemberRoles(ctx, event.ID, callerID)
		if err != nil {
			return EventAuthContext{}, err
		}
		authCtx.Roles = roles
		authCtx.IsMember = len(roles) > 0 || authCtx.IsOwner
	}

	if gameID > 0 && callerID != 0 {
		seated, err := s.gameService.IsPlayerInGame(ctx, gameID, callerID)
		if err != nil {
			return EventAuthContext{}, err
		}
		authCtx.IsSeatedInGame = seated
	}

	return authCtx, nil
}

func (authCtx EventAuthContext) isStaff() bool {
	return slices.Contains(authCtx.Roles, models.EventReferee) ||
		slices.Contains(authCtx.Roles, models.EventObserver)
}

func (authCtx EventAuthContext) isParticipating() bool {
	return authCtx.IsMember || authCtx.IsOwner
}

func (authCtx EventAuthContext) hasRole(role models.EventRole) bool {
	return slices.Contains(authCtx.Roles, role)
}

func (authCtx EventAuthContext) canAccessEvent() bool {
	if authCtx.Event.Visibility == models.EventVisibilityPublic {
		return true
	}
	return authCtx.isParticipating()
}

func (authCtx EventAuthContext) hasCapability(cap models.Capability) bool {
	if !authCtx.canAccessEvent() {
		return false
	}

	switch cap {
	case models.CapManageEvent, models.CapManageMembers:
		return authCtx.IsOwner
	case models.CapCreateGame:
		if authCtx.Event.FinishedAt != nil {
			return false
		}
		if !slices.Contains(models.AllowsNewGameInEvent, authCtx.Event.Type) {
			return false
		}
		return authCtx.hasRole(models.EventPlayer)
	case models.CapAbortGame:
		if !slices.Contains(models.AllowsAbortGameInEvent, authCtx.Event.Type) {
			return false
		}
		return authCtx.hasRole(models.EventReferee)
	case models.CapFinishGame:
		if !slices.Contains(models.AllowsFinishGameInEvent, authCtx.Event.Type) {
			return false
		}
		return authCtx.IsSeatedInGame || authCtx.hasRole(models.EventReferee)
	case models.CapUndoRound:
		return authCtx.hasRole(models.EventReferee)
	case models.CapWatchLive, models.CapWatchResults:
		return authCtx.canViewContent(authCtx.Event.GamesReach, authCtx.baseWatchGamesAccess())
	case models.CapViewRatings:
		return authCtx.canViewContent(authCtx.Event.RatingsReach, authCtx.baseViewRatingsAccess())
	default:
		return false
	}
}

func (authCtx EventAuthContext) baseWatchGamesAccess() bool {
	switch authCtx.Event.Type {
	case models.EventRegular:
		return authCtx.isParticipating()
	case models.EventTournament:
		return false
	default:
		return false
	}
}

func (authCtx EventAuthContext) baseViewRatingsAccess() bool {
	switch authCtx.Event.Type {
	case models.EventRegular:
		return authCtx.isParticipating()
	case models.EventTournament:
		return authCtx.isParticipating()
	default:
		return false
	}
}

func (authCtx EventAuthContext) canViewContent(reach models.ContentReach, isAllowedForMembers bool) bool {
	if authCtx.isStaff() {
		return true
	}

	switch reach {
	case models.ContentReachStaff:
		return false
	case models.ContentReachMembers:
		return isAllowedForMembers
	case models.ContentReachPublic:
		if authCtx.Event.Visibility == models.EventVisibilityPublic {
			return true
		}
		return isAllowedForMembers
	default:
		return false
	}
}
