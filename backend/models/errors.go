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

	ErrWrongClaims        = errors.New("token claims are of wrong type")
	ErrWrongSigningMethod = errors.New("invalid signing method")
	ErrUserIDNotFound     = errors.New("user id not found in context")
	ErrUserIDInvalidType  = errors.New("user id in context is of invalid type")
	ErrEmptyToken         = errors.New("token is empty")

	ErrJWTKeyEmpty     = errors.New("jwt key is empty")
	ErrGeneratingToken = errors.New("error generating token")
)

var Errors = map[error]int{
	// 400
	ErrBadCredentials:   http.StatusBadRequest,
	ErrUserNotInPlayers: http.StatusBadRequest,
	ErrMinScore:         http.StatusBadRequest,

	// 403
	ErrEmptyToken: http.StatusUnauthorized,

	// 404
	ErrGameNotFound:   http.StatusNotFound,
	ErrPlayerNotFound: http.StatusNotFound,
	ErrGameNotLive:    http.StatusNotFound,

	// 409
	ErrPlayerExists: http.StatusConflict,
}
