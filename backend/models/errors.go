package models

import (
	"errors"
	"fmt"
	"net/http"
)

// public errors. Their texts go to the server log; the client gets the code
// from Errors and words the message itself.
var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrLoginRequired  = errors.New("login required")

	ErrGameNotFound             = errors.New("game not found")
	ErrGameNotLive              = errors.New("game is not live")
	ErrGameNotEnded             = errors.New("game is not ended")
	ErrValidateEnterRoundResult = errors.New("round result validate")
	ErrValidatePenalty          = errors.New("penalty validate")

	ErrPlayerNotFound   = errors.New("player not found")
	ErrBadCredentials   = errors.New("email or password is incorrect")
	ErrUserNotInPlayers = errors.New("cannot create game without yourself")
	// ErrPlayerExists means the email already has an account. Never shown to
	// the client: SignUp answers as for a new email, so the response doesn't
	// reveal which emails are registered.
	ErrPlayerExists     = errors.New("player with this email already exists")
	ErrEmptyName        = errors.New("name cannot be empty")
	ErrLongName         = errors.New("name cannot exceed 50 characters")
	ErrEmptyEmail       = errors.New("email cannot be empty")
	ErrLongEmail        = errors.New("email cannot exceed 50 characters")
	ErrInvalidEmail     = errors.New("email is not a valid address")
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrLongPassword     = errors.New("password cannot exceed 72 bytes")
	ErrEmailNotVerified = errors.New("email not verified")
	ErrInvalidAuthToken = errors.New("invalid or expired token")
	// ErrAuthMailLimit means an auth email wasn't sent because the address has
	// had too many lately (models.AuthMailLimits). Never shown to the client:
	// the response says "sent" anyway.
	ErrAuthMailLimit  = errors.New("too many auth emails to this address")
	ErrSessionExpired = errors.New("session expired, please log in again")
	ErrMailFromEmpty  = errors.New("mail_from is empty")

	ErrJWTKeyEmpty    = errors.New("jwt key is empty")
	ErrSiteURLInvalid = errors.New("site_url must be like https://macra.pro (scheme and host, no path)")

	// ErrLogChanged means a round entry or an undo was made from an outdated
	// table or log: someone else has entered or undone an entry since.
	ErrLogChanged = errors.New("the game has changed since you opened it, check the log and try again")
	// ErrPlayerInAnotherGame means a game can't become live (a new game, or
	// an ended one by undoing its final round) because one of its players
	// already has a live game in the event. Sent with the player's name.
	ErrPlayerInAnotherGame = errors.New("player is in another live game")

	// ErrValidatePlayers means the seating isn't four distinct players.
	ErrValidatePlayers = errors.New("validate game players")

	ErrEventNotFound        = errors.New("event not found")
	ErrEventForbidden       = errors.New("forbidden event")
	ErrEventAlreadyFinished = errors.New("event already finished")
	ErrEventHasOngoingGames = errors.New("event has ongoing games")
	// ErrMemberInCurrentGame is sent with the member's name.
	ErrMemberInCurrentGame = errors.New("cannot remove member who plays game right now")
	// ErrMissingPlayerRole is sent with the player's name.
	ErrMissingPlayerRole   = errors.New("missing 'player' role")
	ErrGameNotInEvent      = errors.New("game not in event")
	ErrInvalidEventType    = errors.New("invalid event type")
	ErrInvalidContentReach = errors.New("invalid content reach")
	ErrInvalidVisibility   = errors.New("invalid event visibility")
	ErrLongEventName       = errors.New("event name too long")
	ErrLongDescription     = errors.New("event description too long")
	ErrInvalidGameStatus   = errors.New("invalid game status")

	ErrNoRounds = errors.New("no rounds to undo")
	// ErrGameHasRounds means a seated player can't abort the game any more: a
	// round has been entered (referees still can).
	ErrGameHasRounds = errors.New("game has rounds")
	ErrTooLate       = errors.New("too late to undo")
)

// APIError is how a public error reaches the client: the HTTP status and a
// stable code, sent as {"error": "<code>"} along with the error's data
// (WithData). The frontend maps codes to texts.
type APIError struct {
	Status int
	Code   string
}

// ErrInternal is the APIError for any error not in Errors. Its text stays in
// the server log.
var ErrInternal = APIError{Status: http.StatusInternalServerError, Code: "internal"}

var Errors = map[error]APIError{
	// 400
	ErrInvalidRequest:           {http.StatusBadRequest, "invalid_request"},
	ErrBadCredentials:           {http.StatusBadRequest, "bad_credentials"},
	ErrUserNotInPlayers:         {http.StatusBadRequest, "user_not_in_players"},
	ErrValidateEnterRoundResult: {http.StatusBadRequest, "invalid_round"},
	ErrValidatePenalty:          {http.StatusBadRequest, "invalid_penalty"},
	ErrValidatePlayers:          {http.StatusBadRequest, "invalid_players"},
	ErrNoRounds:                 {http.StatusBadRequest, "no_rounds"},
	ErrTooLate:                  {http.StatusBadRequest, "too_late"},
	ErrEmptyName:                {http.StatusBadRequest, "empty_name"},
	ErrLongName:                 {http.StatusBadRequest, "long_name"},
	ErrEmptyEmail:               {http.StatusBadRequest, "empty_email"},
	ErrLongEmail:                {http.StatusBadRequest, "long_email"},
	ErrInvalidEmail:             {http.StatusBadRequest, "invalid_email"},
	ErrEmptyPassword:            {http.StatusBadRequest, "empty_password"},
	ErrLongPassword:             {http.StatusBadRequest, "long_password"},
	ErrMissingPlayerRole:        {http.StatusBadRequest, "missing_player_role"},
	ErrGameNotInEvent:           {http.StatusBadRequest, "game_not_in_event"},
	ErrInvalidEventType:         {http.StatusBadRequest, "invalid_event_type"},
	ErrInvalidContentReach:      {http.StatusBadRequest, "invalid_content_reach"},
	ErrInvalidVisibility:        {http.StatusBadRequest, "invalid_visibility"},
	ErrLongEventName:            {http.StatusBadRequest, "long_event_name"},
	ErrLongDescription:          {http.StatusBadRequest, "long_description"},
	ErrInvalidGameStatus:        {http.StatusBadRequest, "invalid_game_status"},
	ErrInvalidAuthToken:         {http.StatusBadRequest, "invalid_auth_token"},

	// 401
	ErrSessionExpired: {http.StatusUnauthorized, "session_expired"},
	ErrLoginRequired:  {http.StatusUnauthorized, "login_required"},

	// 403
	ErrEventForbidden:   {http.StatusForbidden, "event_forbidden"},
	ErrEmailNotVerified: {http.StatusForbidden, "email_not_verified"},

	// 404
	ErrGameNotFound:   {http.StatusNotFound, "game_not_found"},
	ErrPlayerNotFound: {http.StatusNotFound, "player_not_found"},
	ErrGameNotLive:    {http.StatusNotFound, "game_not_live"},
	ErrGameNotEnded:   {http.StatusNotFound, "game_not_ended"},
	ErrEventNotFound:  {http.StatusNotFound, "event_not_found"},

	// 409
	ErrLogChanged:           {http.StatusConflict, "log_changed"},
	ErrPlayerInAnotherGame:  {http.StatusConflict, "player_in_another_game"},
	ErrMemberInCurrentGame:  {http.StatusConflict, "member_in_current_game"},
	ErrEventAlreadyFinished: {http.StatusConflict, "event_already_finished"},
	ErrEventHasOngoingGames: {http.StatusConflict, "event_has_ongoing_games"},
	ErrGameHasRounds:        {http.StatusConflict, "game_has_rounds"},
}

// ErrorData is sent to the client next to an error's code, e.g. a player's
// name for the message.
type ErrorData map[string]any

// DataError is an error with data for the client. errors.Is and errors.As see
// the error it wraps.
type DataError struct {
	Err  error
	Data ErrorData
}

// WithData returns err with data for the client, e.g.
// WithData(ErrPlayerInAnotherGame, ErrorData{"name": name}).
func WithData(err error, data ErrorData) error {
	return &DataError{Err: err, Data: data}
}

func (e *DataError) Error() string {
	return fmt.Sprintf("%v %v", e.Err, map[string]any(e.Data))
}

func (e *DataError) Unwrap() error {
	return e.Err
}
