package models

import (
	"errors"
	"net/http"
)

// public errors.
var (
	ErrGameNotFound = errors.New("game not found")
	ErrGameNotLive  = errors.New("game is not live")
	ErrMinScore     = errors.New("hand score should be 8 or higher")

	ErrPlayerNotFound   = errors.New("player not found")
	ErrBadCredentials   = errors.New("email or password is incorrect")
	ErrUserNotInPlayers = errors.New("cannot create game without yourself")
	ErrPlayerExists     = errors.New("player with this email already exists")

	ErrJWTKeyEmpty = errors.New("jwt key is empty")

	ErrOutcomeExists = errors.New("non-penalty outcome already exists")

	ErrValidatePlayers = errors.New("validate game players")
)

var Errors = map[error]int{
	// 400
	ErrBadCredentials:   http.StatusBadRequest,
	ErrUserNotInPlayers: http.StatusBadRequest,
	ErrMinScore:         http.StatusBadRequest,
	ErrValidatePlayers:  http.StatusBadRequest,

	// 404
	ErrGameNotFound:   http.StatusNotFound,
	ErrPlayerNotFound: http.StatusNotFound,
	ErrGameNotLive:    http.StatusNotFound,

	// 409
	ErrPlayerExists:  http.StatusConflict,
	ErrOutcomeExists: http.StatusConflict,
}
