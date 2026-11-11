package models

type EventRole string

var (
	EventPlayer   EventRole = "player"
	EventReferee  EventRole = "referee"
	EventObserver EventRole = "observer"
)

type EventMember struct {
	UserID int         `json:"user_id" db:"user_id"`
	Name   string      `json:"name" db:"name"`
	Roles  []EventRole `json:"roles" db:"roles"`
}
