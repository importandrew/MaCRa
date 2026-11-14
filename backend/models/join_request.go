package models

import "time"

// JoinRequestStatus is a stored request's status. An accepted request isn't
// stored: the user is a player then.
type JoinRequestStatus string

const (
	JoinRequestPending  JoinRequestStatus = "pending"
	JoinRequestDeclined JoinRequestStatus = "declined"
)

// JoinRequestState is what a user may do about joining an event by a
// request (EventMe.JoinRequest).
type JoinRequestState string

const (
	// JoinStateAvailable: the user may ask to join (a visitor signs in
	// first).
	JoinStateAvailable JoinRequestState = "available"
	// JoinStatePending: the user's request waits for the owner's answer;
	// the user may withdraw it.
	JoinStatePending JoinRequestState = "pending"
	// JoinStateDeclined: the owner declined the user's request.
	JoinStateDeclined JoinRequestState = "declined"
	// JoinStateFull: requests are on, but the event has no places left.
	JoinStateFull JoinRequestState = "full"
)

// JoinRequest is a pending request to join, as the owner sees it.
type JoinRequest struct {
	UserID    int       `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// JoiningSettings are how an event takes new players besides the owner
// adding them by hand: requests to join and the player limit (Manage →
// Members).
type JoiningSettings struct {
	AcceptRequests bool `json:"accept_requests"`
	// PlayerLimit is nil for no limit.
	PlayerLimit *int `json:"player_limit"`
}

// Validate returns ErrInvalidPlayerLimit for a limit below 1.
func (s JoiningSettings) Validate() error {
	if s.PlayerLimit != nil && *s.PlayerLimit < 1 {
		return ErrInvalidPlayerLimit
	}
	return nil
}
