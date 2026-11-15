package event

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/service/game"
)

// EventAuthContext is what the capability rules know about the caller in an
// event, and about a game when the check is about one (Game != nil).
type EventAuthContext struct {
	Event    models.Event
	CallerID int
	Roles    []models.EventRole
	IsOwner  bool
	IsMember bool
	// IsSeatedInGame: the caller sits at the game now (their own seat, or
	// one they took as a substitute): they enter its rounds.
	IsSeatedInGame bool
	// HasPlaceInGame: the game is the caller's: their own seat (also while
	// a substitute sits there) or one they sat at. They see it.
	HasPlaceInGame bool
	// FinishedInTour: the caller's own game in the live tour has ended (they
	// may watch its other tables); set for a check about a game of that tour.
	FinishedInTour bool
	// Game is the game the capability is about, nil for an event-level check.
	Game *models.GameSummary
	// Now is the time the rules are checked at (the undo time limit).
	Now time.Time
}

func (s *Service) buildAuthContext(ctx context.Context, event models.Event, callerID, gameID int) (EventAuthContext, error) {
	var roles []models.EventRole
	if callerID != 0 {
		var err error
		if roles, err = s.repo.GetMemberRoles(ctx, event.ID, callerID); err != nil {
			return EventAuthContext{}, err
		}
	}
	authCtx := eventAuthContext(event, callerID, roles)

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
		authCtx.HasPlaceInGame = seated
		if !seated {
			if authCtx.HasPlaceInGame, err = s.gameService.HasPlace(ctx, gameID, callerID); err != nil {
				return EventAuthContext{}, err
			}
		}
	}

	if authCtx.Game != nil && authCtx.Game.TourStatus != nil && *authCtx.Game.TourStatus == models.TourLive &&
		!authCtx.HasPlaceInGame && callerID != 0 {
		tour, err := s.repo.LiveTour(ctx, event.ID, callerID)
		switch {
		case errors.Is(err, models.ErrTourNotFound):
		case err != nil:
			return EventAuthContext{}, err
		default:
			authCtx.FinishedInTour = tour.Finished
		}
	}

	return authCtx, nil
}

// eventAuthContext is the event-level auth context of a caller whose roles in
// the event are known (none for a visitor, callerID 0).
func eventAuthContext(event models.Event, callerID int, roles []models.EventRole) EventAuthContext {
	authCtx := EventAuthContext{
		Event:    event,
		CallerID: callerID,
		IsOwner:  callerID != 0 && event.OwnerUserID == callerID,
		Now:      time.Now(),
	}
	if callerID != 0 {
		authCtx.Roles = roles
		authCtx.IsMember = len(roles) > 0 || authCtx.IsOwner
	}
	return authCtx
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
	case g.TourStatus != nil:
		// A tournament game: while its tour is live, also after the game
		// ended; never once the tour is finished (area 8).
		return *g.TourStatus == models.TourLive &&
			(g.Status == models.GameLive || g.Status == models.GameEnded)
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
		models.CapEnterRound, models.CapEnterPenalty, models.CapRefereeEntry, models.CapViewGame:
		return authCtx.hasGameCapability(capability)
	case models.CapWatchLive, models.CapWatchResults:
		// Whoever can open the event sees its games, but a tour's other
		// tables are hidden while it is played.
		return !authCtx.tourTableHidden()
	case models.CapViewTournament:
		return authCtx.Event.Type == models.EventTournament && authCtx.isStaff()
	case models.CapRunTournament:
		// Running it takes the referee role, also for the owner (area 1).
		return authCtx.Event.Type == models.EventTournament &&
			authCtx.Event.FinishedAt == nil && authCtx.hasRole(models.EventReferee)
	case models.CapViewRatings:
		// Whoever can open the event sees its rating.
		return true
	default:
		return false
	}
}

// hasGameCapability checks a capability about the game in authCtx.Game.
func (authCtx EventAuthContext) hasGameCapability(capability models.Capability) bool {
	switch capability {
	case models.CapAbortGame:
		return authCtx.canAbortGame()
	case models.CapFinishGame:
		return authCtx.canFinishGame()
	case models.CapUndoRound:
		return authCtx.canUndoRound()
	case models.CapEnterRound:
		return authCtx.IsSeatedInGame && authCtx.isGameLive()
	case models.CapEnterPenalty:
		return authCtx.canEnterPenalty()
	case models.CapRefereeEntry:
		// A referee enters from the game's view into any live game of the
		// event, also one they play in.
		return authCtx.hasRole(models.EventReferee) && authCtx.isGameLive()
	case models.CapViewGame:
		return authCtx.canViewGame()
	default:
		return false
	}
}

// canAbortGame: referees at any time; the seated players as the event's
// settings allow.
func (authCtx EventAuthContext) canAbortGame() bool {
	if !slices.Contains(models.AllowsAbortGameInEvent, authCtx.Event.Type) || !authCtx.isGameLive() {
		return false
	}
	return authCtx.hasRole(models.EventReferee) ||
		(authCtx.IsSeatedInGame && authCtx.seatedMayAbort(authCtx.Game.LogSize))
}

// canFinishGame: referees at any time; the seated players unless the event's
// settings leave it to referees.
func (authCtx EventAuthContext) canFinishGame() bool {
	if !authCtx.isGameLive() {
		return false
	}
	// In a tournament nobody finishes a game early; a game stopped at its
	// tour's time is finished by its table or a referee.
	if authCtx.Event.Type == models.EventTournament {
		return authCtx.Game.StoppedAtTime &&
			(authCtx.IsSeatedInGame || authCtx.hasRole(models.EventReferee))
	}
	if !slices.Contains(models.AllowsFinishGameInEvent, authCtx.Event.Type) {
		return false
	}
	return authCtx.hasRole(models.EventReferee) ||
		(authCtx.IsSeatedInGame && authCtx.Event.Settings.FinishGame == models.FinishBySeated)
}

// canUndoRound: referees, also the final round of a game that has just
// ended; the seated players while it is live, as the event's settings allow.
func (authCtx EventAuthContext) canUndoRound() bool {
	return (authCtx.hasRole(models.EventReferee) && authCtx.canUndoInGame()) ||
		(authCtx.IsSeatedInGame && authCtx.isGameLive() &&
			authCtx.Event.Settings.UndoRound == models.UndoBySeated)
}

// canEnterPenalty: at one's own table, the seated players unless the event's
// settings leave penalties to referees, and a seated referee. (A referee
// enters into any game from its view: models.CapRefereeEntry.)
func (authCtx EventAuthContext) canEnterPenalty() bool {
	return authCtx.isGameLive() && authCtx.IsSeatedInGame &&
		(authCtx.hasRole(models.EventReferee) || authCtx.Event.Settings.Penalties == models.PenaltiesBySeated)
}

// canViewGame: whoever can open the event sees its games (a seated player
// their own one in any case), but for a tour's other tables while it is
// played (tourTableHidden).
func (authCtx EventAuthContext) canViewGame() bool {
	return authCtx.Game != nil && !authCtx.tourTableHidden()
}

// tourTableHidden reports whether the game is a tournament table the caller
// may not see now (docs/user-scenarios.md, area 8, "Visibility"): while a
// tour is live, its live tables are seen only by the staff, the table
// itself and the players whose own game in the tour has ended; a drafted
// seating only by the staff, its tables included. Whatever the event's
// visibility, so signing out doesn't get around it.
func (authCtx EventAuthContext) tourTableHidden() bool {
	g := authCtx.Game
	if g == nil || g.TourStatus == nil || authCtx.isStaff() {
		return false
	}
	switch *g.TourStatus {
	case models.TourDraft:
		return true
	case models.TourLive:
		return g.Status == models.GameLive && !authCtx.HasPlaceInGame && !authCtx.FinishedInTour
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
