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

	ErrPlayerNotFound   = errors.New("player not found")
	ErrBadCredentials   = errors.New("email or password is incorrect")
	ErrUserNotInPlayers = errors.New("cannot create game without yourself")
	ErrPlayerExists     = errors.New("player with this email already exists")
	ErrEmptyName        = errors.New("name cannot be empty")
	ErrLongName         = errors.New("name cannot exceed 50 characters")

	ErrJWTKeyEmpty = errors.New("jwt key is empty")

	ErrOutcomeExists = errors.New("non-penalty outcome already exists")

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

	// 409
	ErrPlayerExists:         http.StatusConflict,
	ErrOutcomeExists:        http.StatusConflict,
	ErrMemberInCurrentGame:  http.StatusConflict,
	ErrEventAlreadyFinished: http.StatusConflict,
	ErrEventHasOngoingGames: http.StatusConflict,
}
