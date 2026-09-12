package models

import (
	"time"
)

const (
	EventIDLockedDefault = 0
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

type Event struct {
	ID          int             `db:"id" json:"id"`
	Type        EventType       `db:"type" json:"type"`
	Name        string          `db:"name" json:"name"`
	Description string          `db:"description" json:"description,omitempty"`
	OwnerUserID int             `db:"owner_user_id" json:"owner_user_id"`
	Visibility   EventVisibility `db:"visibility" json:"visibility"`
	RatingsReach ContentReach    `db:"ratings_reach" json:"ratings_reach"`
	GamesReach   ContentReach    `db:"games_reach" json:"games_reach"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	FinishedAt  *time.Time      `db:"finished_at" json:"finished_at"`
}

type EventListItem struct {
	Event
	IsParticipating bool        `json:"is_participating"`
	IsOwner         bool        `json:"is_owner"`
	Roles           []EventRole `json:"roles"`
}

type EventMe struct {
	IsOwner      bool         `json:"is_owner"`
	Roles        []EventRole  `json:"roles"`
	Capabilities []Capability `json:"capabilities"`
}

type EventCreateParams struct {
	OwnerID     int
	Type        EventType
	Name        string
	Description string
	Visibility   EventVisibility
	RatingsReach *ContentReach
	GamesReach   *ContentReach
}

type UpdateEventParams struct {
	ID           int
	Name         *string
	Description  *string
	Visibility   *EventVisibility
	RatingsReach *ContentReach
	GamesReach   *ContentReach
}

// EventRatingEntry is one row in the event standings table.
type EventRatingEntry struct {
	Rank     int     `json:"rank"`
	PlayerID int     `json:"player_id"`
	Name     string  `json:"name"`
	Games    int     `json:"games"`
	TP       float32 `json:"tp"`
	GP       int     `json:"gp"`
}

// EventPlayerResultRow is a per-game result used to build event ratings.
type EventPlayerResultRow struct {
	PlayerID int     `db:"player_id"`
	Name     string  `db:"name"`
	TP       float32 `db:"tp"`
	GP       int     `db:"gp"`
}
