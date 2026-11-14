package models

type EventRole string

var (
	EventPlayer   EventRole = "player"
	EventReferee  EventRole = "referee"
	EventObserver EventRole = "observer"
)

// ValidEventRoles is every member role.
var ValidEventRoles = []EventRole{EventPlayer, EventReferee, EventObserver}

type EventMember struct {
	UserID int         `json:"user_id" db:"user_id"`
	Name   string      `json:"name" db:"name"`
	Roles  []EventRole `json:"roles" db:"roles"`
}
