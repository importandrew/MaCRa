package models

import (
	"errors"
	"net/http"
)

// public errors.
var (
	ErrGameNotFound             = errors.New("game not found")
	ErrGameNotLive              = errors.New("game is not live")
	ErrGameNotEnded             = errors.New("game is not ended")
	ErrValidateEnterRoundResult = errors.New("round result validate")
	ErrValidatePenalty          = errors.New("penalty validate")

	ErrPlayerNotFound   = errors.New("player not found")
	ErrBadCredentials   = errors.New("email or password is incorrect")
	ErrUserNotInPlayers = errors.New("cannot create game without yourself")
	ErrPlayerExists     = errors.New("player with this email already exists")
	ErrEmptyName        = errors.New("name cannot be empty")
	ErrLongName         = errors.New("name cannot exceed 50 characters")
	ErrEmailNotVerified = errors.New("email not verified")
	ErrInvalidOrigin    = errors.New("invalid origin")
	ErrInvalidAuthToken = errors.New("invalid or expired token")
	ErrSessionExpired   = errors.New("session expired, please log in again")
	ErrMailFromEmpty    = errors.New("mail_from is empty")

	ErrJWTKeyEmpty = errors.New("jwt key is empty")

	// ErrLogChanged means a round entry or an undo was made from an outdated
	// table or log: someone else has entered or undone an entry since.
	ErrLogChanged = errors.New("the game has changed since you opened it, check the log and try again")
	// ErrUndoFinalRound means the final round can't be undone because the
	// game can't become live again (a player has started another game).
	ErrUndoFinalRound = errors.New("cannot undo the final round")

	ErrValidatePlayers = errors.New("validate game players")

	ErrEventNotFound         = errors.New("event not found")
	ErrEmptyEvent            = errors.New("empty event")
	ErrEventForbidden        = errors.New("forbidden event")
	ErrImmutableEvent        = errors.New("event cannot be modified")
	ErrEventAlreadyFinished  = errors.New("event already finished")
	ErrEventHasOngoingGames  = errors.New("event has ongoing games")
	ErrEventMembersForbidden = errors.New("only the event owner can manage members")
	ErrMemberInCurrentGame   = errors.New("cannot remove member who plays game right now")
	ErrMissingPlayerRole     = errors.New("missing 'player' role")
	ErrPlayerNotInEvent      = errors.New("player not in event")
	ErrGameNotInEvent        = errors.New("game not in event")
	ErrInvalidEventType      = errors.New("invalid event type")
	ErrInvalidContentReach   = errors.New("invalid content reach")
	ErrInvalidGameStatus     = errors.New("invalid game status")

	ErrNoRounds = errors.New("no rounds to undo")
	ErrTooLate  = errors.New("too late to undo")
)

var Errors = map[error]int{
	// 400
	ErrBadCredentials:           http.StatusBadRequest,
	ErrUserNotInPlayers:         http.StatusBadRequest,
	ErrValidateEnterRoundResult: http.StatusBadRequest,
	ErrValidatePenalty:          http.StatusBadRequest,
	ErrValidatePlayers:          http.StatusBadRequest,
	ErrNoRounds:                 http.StatusBadRequest,
	ErrTooLate:                  http.StatusBadRequest,
	ErrEmptyName:                http.StatusBadRequest,
	ErrLongName:                 http.StatusBadRequest,
	ErrMissingPlayerRole:        http.StatusBadRequest,
	ErrPlayerNotInEvent:         http.StatusBadRequest,
	ErrGameNotInEvent:           http.StatusBadRequest,
	ErrInvalidEventType:         http.StatusBadRequest,
	ErrInvalidContentReach:      http.StatusBadRequest,
	ErrInvalidGameStatus:        http.StatusBadRequest,
	ErrInvalidOrigin:            http.StatusBadRequest,
	ErrInvalidAuthToken:         http.StatusBadRequest,

	// 401
	ErrSessionExpired: http.StatusUnauthorized,
	// 404
	ErrGameNotFound:   http.StatusNotFound,
	ErrPlayerNotFound: http.StatusNotFound,
	ErrGameNotLive:    http.StatusNotFound,
	ErrGameNotEnded:   http.StatusNotFound,
	ErrEventNotFound:  http.StatusNotFound,

	// 403
	ErrEventForbidden:        http.StatusForbidden,
	ErrEventMembersForbidden: http.StatusForbidden,
	ErrImmutableEvent:        http.StatusForbidden,
	ErrEmailNotVerified:      http.StatusForbidden,

	// 409
	ErrPlayerExists:         http.StatusConflict,
	ErrLogChanged:           http.StatusConflict,
	ErrUndoFinalRound:       http.StatusConflict,
	ErrMemberInCurrentGame:  http.StatusConflict,
	ErrEventAlreadyFinished: http.StatusConflict,
	ErrEventHasOngoingGames: http.StatusConflict,
}

// PlayerInAnotherGameError is returned when the final round can't be undone
// because one of the game's players has started another game. It is
// ErrUndoFinalRound for errors.Is (and so a 409), and its text names the
// player, for the client to show as is.
type PlayerInAnotherGameError struct {
	Name string
}

func (e PlayerInAnotherGameError) Error() string {
	return e.Name + " has already started another game"
}

// Is reports whether target is ErrUndoFinalRound.
func (e PlayerInAnotherGameError) Is(target error) bool {
	return target == ErrUndoFinalRound
}
