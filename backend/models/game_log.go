package models

import (
	"time"
)

type GameLog struct {
	ID        int            `json:"id"`
	Players   []PlayerHeader `json:"players"`
	RoundLogs []RoundLog     `json:"round_logs"`
	// Undone are the entries undone so far, in the order they were undone;
	// the client shows each one in the log at the time of its undo.
	Undone []UndoneEntry `json:"undone"`
}

// UndoneEntry is a round log entry that has been undone: what it was, who
// entered it and who undid it, and when. Either user is nil when unknown
// (an entry from before authors were recorded, a deleted account).
type UndoneEntry struct {
	Round          int           `json:"round"`
	Outcome        Outcome       `json:"outcome"`
	Winner         int           `json:"winner"`
	Loser          int           `json:"loser"`
	Value          int           `json:"value"`
	ForceNextRound bool          `json:"force_next_round"`
	EnteredBy      *PlayerHeader `json:"entered_by"`
	EnteredAt      time.Time     `json:"entered_at"`
	UndoneBy       *PlayerHeader `json:"undone_by"`
	UndoneAt       time.Time     `json:"undone_at"`
	// ScoresBefore and ScoresAfter are the players' scores before and after
	// the entry, by starting seat (east, south, west, north).
	ScoresBefore [4]int `json:"scores_before"`
	ScoresAfter  [4]int `json:"scores_after"`
}

type PlayerHeader struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RoundLog struct {
	ID             int       `json:"id"`
	Round          int       `json:"round"`
	EventSequence  int       `json:"event_sequence"`
	Outcome        Outcome   `json:"outcome"`
	Winner         int       `json:"winner"`
	Loser          int       `json:"loser"`
	Value          int       `json:"value"`
	ForceNextRound bool      `json:"force_next_round"`
	EastScore      int       `json:"east_score"`
	SouthScore     int       `json:"south_score"`
	WestScore      int       `json:"west_score"`
	NorthScore     int       `json:"north_score"`
	CreatedAt      time.Time `json:"created_at"`
	// EnteredBy is who entered the entry; nil for an entry saved before
	// authors were recorded (or by a deleted account).
	EnteredBy *PlayerHeader `json:"entered_by"`
}

func (l RoundLog) Score(wind Wind) int {
	switch wind {
	case East:
		return l.EastScore
	case South:
		return l.SouthScore
	case West:
		return l.WestScore
	case North:
		return l.NorthScore
	default:
		return -1234
	}
}
