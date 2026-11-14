package event

import (
	"context"
	"slices"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/service/game"
)

// EventAuthContext is what the capability rules know about the caller in an
// event, and about a game when the check is about one (Game != nil).
type EventAuthContext struct {
	Event          models.Event
	CallerID       int
	Roles          []models.EventRole
	IsOwner        bool
	IsMember       bool
	IsSeatedInGame bool
	// Game is the game the capability is about, nil for an event-level check.
	Game *models.GameSummary
	// Now is the time the rules are checked at (the undo time limit).
	Now time.Time
}

func (s *Service) buildAuthContext(ctx context.Context, event models.Event, callerID, gameID int) (EventAuthContext, error) {
	authCtx := EventAuthContext{
		Event:    event,
		CallerID: callerID,
		IsOwner:  event.OwnerUserID == callerID,
		Now:      time.Now(),
	}

	if callerID != 0 {
		roles, err := s.repo.GetMemberRoles(ctx, event.ID, callerID)
		if err != nil {
			return EventAuthContext{}, err
		}
		authCtx.Roles = roles
		authCtx.IsMember = len(roles) > 0 || authCtx.IsOwner
	}

	if gameID > 0 {
		summary, err := s.gameService.Summary(ctx, gameID)
		if err != nil {
			return EventAuthContext{}, err
		}
		authCtx.Game = &summary
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

// isStaff is about seeing: referees, observers and the owner see everything
// in the event, whatever the reaches. Acting as a referee (undo, abort) takes
// the referee role, also for the owner.
func (authCtx EventAuthContext) isStaff() bool {
	return authCtx.IsOwner ||
		slices.Contains(authCtx.Roles, models.EventReferee) ||
		slices.Contains(authCtx.Roles, models.EventObserver)
}

// isGameLive reports whether the check is about a live game.
func (authCtx EventAuthContext) isGameLive() bool {
	return authCtx.Game != nil && authCtx.Game.Status == models.GameLive
}

// canUndoInGame: a live game, or one that ended less than
// models.UndoFinalRoundTimeout ago (it becomes live again).
func (authCtx EventAuthContext) canUndoInGame() bool {
	g := authCtx.Game
	switch {
	case g == nil:
		return false
	case g.Status == models.GameLive:
		return true
	case g.Status == models.GameEnded && g.FinishedAt != nil:
		return g.FinishedAt.After(game.UndoFinalRoundCutoff(authCtx.Now))
	default:
		return false
	}
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

func (authCtx EventAuthContext) hasCapability(capability models.Capability) bool {
	if !authCtx.canAccessEvent() {
		return false
	}

	switch capability {
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
	case models.CapAbortGame, models.CapFinishGame, models.CapUndoRound,
		models.CapEnterRound, models.CapViewGame:
		return authCtx.hasGameCapability(capability)
	case models.CapWatchLive, models.CapWatchResults:
		// A seated player always sees their own game, whatever the games reach.
		if authCtx.IsSeatedInGame {
			return true
		}
		return authCtx.canViewContent(authCtx.Event.GamesReach, authCtx.baseWatchGamesAccess())
	case models.CapViewRatings:
		return authCtx.canViewContent(authCtx.Event.RatingsReach, authCtx.baseViewRatingsAccess())
	default:
		return false
	}
}

// hasGameCapability checks a capability about the game in authCtx.Game.
func (authCtx EventAuthContext) hasGameCapability(capability models.Capability) bool {
	switch capability {
	case models.CapAbortGame:
		// Referees at any time; the seated players as the event's settings
		// allow.
		if !slices.Contains(models.AllowsAbortGameInEvent, authCtx.Event.Type) || !authCtx.isGameLive() {
			return false
		}
		return authCtx.hasRole(models.EventReferee) ||
			(authCtx.IsSeatedInGame && authCtx.seatedMayAbort(authCtx.Game.LogSize))
	case models.CapFinishGame:
		// Referees at any time; the seated players unless the event's
		// settings leave it to referees.
		if !slices.Contains(models.AllowsFinishGameInEvent, authCtx.Event.Type) || !authCtx.isGameLive() {
			return false
		}
		return authCtx.hasRole(models.EventReferee) ||
			(authCtx.IsSeatedInGame && authCtx.Event.Settings.FinishGame == models.FinishBySeated)
	case models.CapUndoRound:
		// Referees, also the final round of a game that has just ended; the
		// seated players while it is live, as the event's settings allow.
		return (authCtx.hasRole(models.EventReferee) && authCtx.canUndoInGame()) ||
			(authCtx.IsSeatedInGame && authCtx.isGameLive() &&
				authCtx.Event.Settings.UndoRound == models.UndoBySeated)
	case models.CapEnterRound:
		return authCtx.IsSeatedInGame && authCtx.isGameLive()
	case models.CapViewGame:
		// A seated player always sees their own game; others by the games
		// reach: watch_live while it is live, watch_results once it isn't.
		if authCtx.Game == nil {
			return false
		}
		if authCtx.IsSeatedInGame {
			return true
		}
		return authCtx.canViewContent(authCtx.Event.GamesReach, authCtx.baseWatchGamesAccess())
	default:
		return false
	}
}

// seatedMayAbort reports whether the event's settings let a seated player
// abort a game with logSize log entries.
func (authCtx EventAuthContext) seatedMayAbort(logSize int) bool {
	switch authCtx.Event.Settings.AbortGame {
	case models.AbortBySeated:
		return true
	case models.AbortBySeatedBeforeFirstRound:
		return logSize == 0
	default:
		return false
	}
}

// seatedAbortOnlyWithoutRounds reports whether a seated player's abort must
// be refused once the game has a round: it is checked again in the same
// statement as the abort (repo/game.Abort).
func (authCtx EventAuthContext) seatedAbortOnlyWithoutRounds() bool {
	return authCtx.Event.Settings.AbortGame == models.AbortBySeatedBeforeFirstRound
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
