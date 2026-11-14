package models

import "time"

// TourStatus is where a tour is: its seating published (games created),
// live (started), finished, or aborted (it doesn't count; the number is
// played again).
type TourStatus string

const (
	TourPublished TourStatus = "published"
	TourLive      TourStatus = "live"
	TourFinished  TourStatus = "finished"
	TourAborted   TourStatus = "aborted"
)

// SeatingMethod is how a tour's tables are made.
type SeatingMethod string

const (
	// SeatingRandom: the active participants shuffled into tables.
	SeatingRandom SeatingMethod = "random"
)

var ValidSeatingMethods = []SeatingMethod{SeatingRandom}

// Tour is one tour of a tournament.
type Tour struct {
	ID          int           `db:"id" json:"id"`
	Number      int           `db:"number" json:"number"`
	Status      TourStatus    `db:"status" json:"status"`
	Seating     SeatingMethod `db:"seating" json:"seating"`
	PublishedAt time.Time     `db:"published_at" json:"published_at"`
	StartedAt   *time.Time    `db:"started_at" json:"started_at"`
	FinishedAt  *time.Time    `db:"finished_at" json:"finished_at"`
}

// SeatsPerTable is the number of players at a table.
const SeatsPerTable = 4

// Seat is a participant at a tour's table.
type Seat struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Number int    `json:"number"`
	// Wind is the starting seat: "east", "south", "west" or "north".
	Wind string `json:"wind"`
}

// TourTable is a table of a tour: its game and the four seats, east first.
type TourTable struct {
	Table  int    `json:"table"`
	GameID int    `json:"game_id"`
	Seats  []Seat `json:"seats"`
}

// TourState is a tournament's tours as its desk shows them: the current tour
// (the last one not aborted; nil before the first) with its tables, and what
// the next tour would be.
type TourState struct {
	Current *Tour       `json:"current"`
	Tables  []TourTable `json:"tables"`
	// Next is the number of the tour to publish next; 0 when all are played.
	Next int `json:"next"`
	// Active is the number of participants who can be seated (not
	// withdrawn); Missing how many players a whole number of tables lacks.
	Active  int `json:"active"`
	Missing int `json:"missing"`
}

// Participant is a seated-to-be participant: a numbered player who hasn't
// withdrawn.
type Participant struct {
	UserID int `db:"user_id"`
	Number int `db:"number"`
}

// TablesFunc makes a tour's tables from the participants to seat: each
// table's four user IDs, east to north.
type TablesFunc func(participants []Participant) ([][4]int, error)

// TourSeat is where a player sits in a tour, for their own table screen.
type TourSeat struct {
	Tour  int `json:"tour"`
	Table int `json:"table"`
}
