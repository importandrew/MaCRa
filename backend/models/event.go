package models

import (
	"time"
)

type EventType string

var (
	EventTournament EventType = "tournament"
	EventRegular    EventType = "regular"
)

var ValidEventTypes = []EventType{EventTournament, EventRegular}
var AllowsNewGameInEvent = []EventType{EventRegular}
var AllowsAbortGameInEvent = []EventType{EventRegular}
var AllowsFinishGameInEvent = []EventType{EventRegular}

type EventVisibility string

const (
	EventVisibilityPrivate EventVisibility = "private"
	EventVisibilityPublic  EventVisibility = "public"
)

var ValidEventVisibilities = []EventVisibility{EventVisibilityPrivate, EventVisibilityPublic}

// Limits of an event's texts, in characters (the frontend mirrors them).
const (
	MaxEventNameLength        = 100
	MaxEventDescriptionLength = 1000
)

type Event struct {
	ID           int             `db:"id" json:"id"`
	Type         EventType       `db:"type" json:"type"`
	Name         string          `db:"name" json:"name"`
	Description  string          `db:"description" json:"description,omitempty"`
	OwnerUserID  int             `db:"owner_user_id" json:"owner_user_id"`
	Visibility   EventVisibility `db:"visibility" json:"visibility"`
	RatingsReach ContentReach    `db:"ratings_reach" json:"ratings_reach"`
	GamesReach   ContentReach    `db:"games_reach" json:"games_reach"`
	Rules        GameRules       `db:"rules" json:"rules"`
	Settings     EventSettings   `db:"settings" json:"settings"`
	// AcceptRequests: the owner takes requests to join (public events only).
	AcceptRequests bool `db:"accept_requests" json:"accept_requests"`
	// PlayerLimit is the most players the event takes by a request or the
	// invite link; nil for no limit.
	PlayerLimit *int       `db:"player_limit" json:"player_limit"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	FinishedAt  *time.Time `db:"finished_at" json:"finished_at"`
}

// EventDetails is an event with what its overview shows besides the event
// itself. PlayerCount is shown to everyone who can open the event, also to
// non-members, who don't get the member list.
type EventDetails struct {
	Event

	PlayerCount int `json:"player_count"`
	// RulesFrozen: the event has games, so its game rules can't change.
	RulesFrozen bool `json:"rules_frozen"`
}

type EventListItem struct {
	Event

	IsParticipating bool        `json:"is_participating"`
	IsOwner         bool        `json:"is_owner"`
	Roles           []EventRole `json:"roles"`
	// LastGameAt is when the last game that wasn't aborted started; nil
	// without games. Lists order active events by it.
	LastGameAt *time.Time `json:"last_game_at"`
}

// EventListRow is an event as the event list reads it: with the caller's
// roles in it and the start of its last game.
type EventListRow struct {
	Event

	Roles      []EventRole `db:"roles"`
	LastGameAt *time.Time  `db:"last_game_at"`
}

type EventMe struct {
	IsOwner      bool         `json:"is_owner"`
	Roles        []EventRole  `json:"roles"`
	Capabilities []Capability `json:"capabilities"`
	// JoinRequest is what the overview offers a non-member about joining by
	// a request; empty (not sent) when it offers nothing.
	JoinRequest JoinRequestState `json:"join_request,omitempty"`
	// PendingRequests is the number of requests to join waiting for an
	// answer, for whoever has CapManageMembers (0, not sent, otherwise).
	PendingRequests int `json:"pending_requests,omitempty"`
}

type EventCreateParams struct {
	OwnerID      int
	Type         EventType
	Name         string
	Description  string
	Visibility   EventVisibility
	RatingsReach *ContentReach
	GamesReach   *ContentReach
	// Rules are the event's game rules; nil for the MCR standard ones.
	Rules *GameRules
}

type UpdateEventParams struct {
	ID           int
	Name         *string
	Description  *string
	Visibility   *EventVisibility
	RatingsReach *ContentReach
	GamesReach   *ContentReach
}

// EventRatingEntry is one row in the event standings table. Rank is 0 (not
// sent) for a player below the event's minimum games.
type EventRatingEntry struct {
	Rank     int     `json:"rank,omitempty"`
	PlayerID int     `json:"player_id"`
	Name     string  `json:"name"`
	Games    int     `json:"games"`
	TP       float32 `json:"tp"`
	GP       int     `json:"gp"`
}

// EventRatings is an event's rating: the players with at least MinGames
// ended games, ranked, and below them, without a place, the players with
// fewer, in the same order, then (for members and staff) the players who
// haven't played yet.
type EventRatings struct {
	Ratings      []EventRatingEntry `json:"ratings"`
	BelowMinimum []EventRatingEntry `json:"below_minimum"`
	MinGames     int                `json:"min_games"`
}

// EventPlayerResultRow is a per-game result used to build event ratings.
type EventPlayerResultRow struct {
	PlayerID int     `db:"player_id"`
	Name     string  `db:"name"`
	TP       float32 `db:"tp"`
	GP       int     `db:"gp"`
}

// EventInvite is what an invite link shows before joining: the event's name
// and type, whether the caller is already a member (false for a visitor),
// and whether the event is full (false for a member).
type EventInvite struct {
	EventID  int       `json:"event_id"`
	Name     string    `json:"name"`
	Type     EventType `json:"type"`
	IsMember bool      `json:"is_member"`
	IsFull   bool      `json:"is_full"`
}
