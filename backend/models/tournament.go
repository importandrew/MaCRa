package models

import (
	"fmt"
	"slices"
	"time"
)

// TournamentSettings are a tournament's own settings, stored as complete
// JSON in events.tournament (nil for a regular event). When each may change
// is in docs/user-scenarios.md, area 8 ("Settings and when they freeze");
// until tours exist, all of them may.
type TournamentSettings struct {
	// Tours is the number of tours (rounds of games).
	Tours int `json:"tours"`
	// TourMinutes is the length of a tour.
	TourMinutes int `json:"tour_minutes"`
	// TimeUp is what happens at a table when the tour's time is up:
	// TimeUpFinishHand, TimeUpExtraHands or TimeUpStopNow.
	TimeUp string `json:"time_up"`
	// ExtraHands is how many hands are played after the one in progress,
	// with TimeUpExtraHands; 0 otherwise.
	ExtraHands int `json:"extra_hands"`
}

// What happens when a tour's time is up.
const (
	// TimeUpFinishHand: the hand in progress is played to the end (the
	// organizer's practice).
	TimeUpFinishHand = "finish_hand"
	// TimeUpExtraHands: the hand in progress and ExtraHands more.
	TimeUpExtraHands = "extra_hands"
	// TimeUpStopNow: the game ends at once, the hand in progress doesn't
	// count (the official rule).
	TimeUpStopNow = "stop_now"
)

var ValidTimeUp = []string{TimeUpFinishHand, TimeUpExtraHands, TimeUpStopNow}

// Limits of the tournament settings (the frontend mirrors them).
const (
	MinTours       = 3
	MaxTours       = 15
	MinTourMinutes = 30
	MaxTourMinutes = 300
	MinExtraHands  = 1
	MaxExtraHands  = 3
)

// The settings a new tournament starts with.
const (
	DefaultTours       = 5
	DefaultTourMinutes = 120
)

// DefaultTournamentSettings returns the settings a new tournament starts
// with: 5 tours of 120 minutes, the hand in progress finished when time is up.
func DefaultTournamentSettings() TournamentSettings {
	return TournamentSettings{
		Tours:       DefaultTours,
		TourMinutes: DefaultTourMinutes,
		TimeUp:      TimeUpFinishHand,
		ExtraHands:  0,
	}
}

// Validate returns an error wrapping ErrInvalidTournamentSettings if a value
// is out of its range or not one of its options.
func (s TournamentSettings) Validate() error {
	if s.Tours < MinTours || s.Tours > MaxTours {
		return fmt.Errorf("%w: tours should be %d to %d, got %d",
			ErrInvalidTournamentSettings, MinTours, MaxTours, s.Tours)
	}
	if s.TourMinutes < MinTourMinutes || s.TourMinutes > MaxTourMinutes {
		return fmt.Errorf("%w: tour_minutes should be %d to %d, got %d",
			ErrInvalidTournamentSettings, MinTourMinutes, MaxTourMinutes, s.TourMinutes)
	}
	if !slices.Contains(ValidTimeUp, s.TimeUp) {
		return fmt.Errorf("%w: time_up should be one of %v, got %q",
			ErrInvalidTournamentSettings, ValidTimeUp, s.TimeUp)
	}
	if s.TimeUp == TimeUpExtraHands {
		if s.ExtraHands < MinExtraHands || s.ExtraHands > MaxExtraHands {
			return fmt.Errorf("%w: extra_hands should be %d to %d, got %d",
				ErrInvalidTournamentSettings, MinExtraHands, MaxExtraHands, s.ExtraHands)
		}
	} else if s.ExtraHands != 0 {
		return fmt.Errorf("%w: extra_hands is for %q only, got %d",
			ErrInvalidTournamentSettings, TimeUpExtraHands, s.ExtraHands)
	}
	return nil
}

// EventStatus is where an event is in its life: a tournament is prepared,
// then in progress from its first tour, then finished; a regular event is in
// progress until it is finished.
type EventStatus string

const (
	EventPreparation EventStatus = "preparation"
	EventInProgress  EventStatus = "in_progress"
	EventFinished    EventStatus = "finished"
)

// eventStatus is the status of an event from its type and times.
func eventStatus(eventType EventType, startedAt, finishedAt *time.Time) EventStatus {
	switch {
	case finishedAt != nil:
		return EventFinished
	case eventType == EventTournament && startedAt == nil:
		return EventPreparation
	default:
		return EventInProgress
	}
}

// TournamentParticipant is a tournament's player as its desk lists them.
type TournamentParticipant struct {
	UserID int    `db:"user_id" json:"user_id"`
	Name   string `db:"name" json:"name"`
	// Number is the participant number from the draw; nil before it, and
	// while the draw is out of date (TournamentParticipants.Drawn).
	Number    *int `db:"number" json:"number"`
	Withdrawn bool `db:"withdrawn" json:"withdrawn"`
}

// TournamentParticipants is a tournament's players with their numbers.
type TournamentParticipants struct {
	// Drawn: the numbers belong to exactly the current players. False before
	// the first draw and after the players changed: the draw is redone.
	Drawn bool `json:"drawn"`
	// Participants by number when drawn, by name otherwise.
	Participants []TournamentParticipant `json:"participants"`
}
