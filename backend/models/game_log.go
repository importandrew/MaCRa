package models

import (
	"time"
)

type GameLog struct {
	ID        int            `json:"id"`
	Players   []PlayerHeader `json:"players"`
	RoundLogs []RoundLog     `json:"round_logs"`
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
}
