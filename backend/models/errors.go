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

	// 404
	ErrGameNotFound:   http.StatusNotFound,
	ErrPlayerNotFound: http.StatusNotFound,
	ErrGameNotLive:    http.StatusNotFound,
	ErrGameNotEnded:   http.StatusNotFound,

	// 409
	ErrPlayerExists:  http.StatusConflict,
	ErrOutcomeExists: http.StatusConflict,
}
