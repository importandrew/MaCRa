package models

import "time"

// TourStatus is where a tour is: its seating generated as a draft (games
// created, seen by the staff only), published (the players see their
// tables), live (started), finished, or aborted (it doesn't count; the number
// is played again).
type TourStatus string

const (
	TourDraft     TourStatus = "draft"
	TourPublished TourStatus = "published"
	TourLive      TourStatus = "live"
	TourFinished  TourStatus = "finished"
	TourAborted   TourStatus = "aborted"
)

// SeatingMethod is how a tour's tables are made.
type SeatingMethod string

const (
	// SeatingRandom: at random, with as few repeat meetings as possible.
	SeatingRandom SeatingMethod = "random"
	// SeatingSwiss: neighbours in the current standings together, avoiding
	// repeat meetings.
	SeatingSwiss SeatingMethod = "swiss"
	// SeatingInterval: by the current standings with a step of the number
	// of tables: table i gets places i, i+k, i+2k, i+3k.
	SeatingInterval SeatingMethod = "interval"
	// SeatingFixed: the tour's tables imported from a file (TourPlan.Tables),
	// winds as there.
	SeatingFixed SeatingMethod = "fixed"
)

var ValidSeatingMethods = []SeatingMethod{SeatingRandom, SeatingSwiss, SeatingInterval, SeatingFixed}

// Seating is how a tour's tables are made: the method, for an interval
// seating its step (nil: the number of tables), and for a Swiss or interval
// seating whether it avoids repeat meetings (a random one always does).
type Seating struct {
	Method       SeatingMethod
	Step         *int
	AvoidRepeats bool
}

// Tour is one tour of a tournament.
type Tour struct {
	ID      int           `db:"id" json:"id"`
	Number  int           `db:"number" json:"number"`
	Status  TourStatus    `db:"status" json:"status"`
	Seating SeatingMethod `db:"seating" json:"seating"`
	// Step is an interval seating's step; nil for the other methods.
	Step *int `db:"seating_step" json:"seating_step"`
	// AvoidRepeats: a Swiss or interval seating moved players apart to
	// avoid repeat meetings.
	AvoidRepeats bool `db:"avoid_repeats" json:"avoid_repeats"`
	// PublishedAt is nil while the seating is a draft.
	PublishedAt *time.Time `db:"published_at" json:"published_at"`
	// Minutes is the tour's length, set when it starts; nil before.
	Minutes    *int       `db:"minutes" json:"minutes"`
	StartedAt  *time.Time `db:"started_at" json:"started_at"`
	FinishedAt *time.Time `db:"finished_at" json:"finished_at"`
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
	// Sitter is who sits at the seat instead of its player (a substitute);
	// nil while the player does.
	Sitter *Player `json:"sitter"`
	// Forfeit: the player gets 0 TP and 0 GP for the game.
	Forfeit bool `json:"forfeit"`
}

// TourTable is a table of a tour: its game and the four seats, east first.
type TourTable struct {
	Table  int `json:"table"`
	GameID int `json:"game_id"`
	// Status is the table's game status: created, live, ended or aborted.
	Status GameStatus `json:"status"`
	// ExtraMinutes is the time a referee added to the table.
	ExtraMinutes int `json:"extra_minutes"`
	// LogSize is the number of the game's log entries, for a referee's
	// change at the table (seating a substitute).
	LogSize int    `json:"log_size"`
	Seats   []Seat `json:"seats"`
	// Repeats counts the table's repeat meetings; RepeatPairs are the pairs
	// who met before, with how many times.
	Repeats     int          `json:"repeats"`
	RepeatPairs []RepeatPair `json:"repeat_pairs"`
	// Places are the lowest and highest places in the standings at the
	// table; nil before any tour is finished.
	Places *PlaceRange `json:"places"`
}

// RepeatPair is two players at a table who have met before, Met times.
type RepeatPair struct {
	A   int `json:"a"`
	B   int `json:"b"`
	Met int `json:"met"`
}

// PlaceRange is the range of places in the standings at a table.
type PlaceRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// TourState is a tournament's tours as its desk shows them: the current tour
// (the last one not aborted; nil before the first) with its tables, and what
// the next tour would be.
type TourState struct {
	Current *Tour       `json:"current"`
	Tables  []TourTable `json:"tables"`
	// Next is the number of the tour to publish next; 0 when all are played.
	Next int `json:"next"`
	// Repeats counts the repeat meetings at the current tour's tables: each
	// earlier meeting of each pair seated together. MaxMet is the most times
	// a pair seated together has met before; PlayersWithRepeats how many
	// players meet someone again.
	Repeats            int `json:"repeats"`
	MaxMet             int `json:"max_met"`
	PlayersWithRepeats int `json:"players_with_repeats"`
	// Active is the number of participants who can be seated (not
	// withdrawn); Missing how many players a whole number of tables lacks.
	Active  int `json:"active"`
	Missing int `json:"missing"`
	// Plan is the tournament's plan: each tour's seating method, 1 to the
	// number of tours. File describes the tables imported from a file; nil
	// without them.
	Plan []TourPlan       `json:"plan"`
	File *SeatingFileInfo `json:"file"`
	// Substitutes are the tournament's substitutes, by name, with the table
	// each sits at now (0: none), for seating them.
	Substitutes []Substitute `json:"substitutes"`
}

// Substitute is a substitute of a tournament and the table they sit at now
// (0 when they sit nowhere).
type Substitute struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Table int    `db:"table_number" json:"table"`
}

// TourPlan is the seating a tour is to have. A tour without a plan of its
// own is random.
type TourPlan struct {
	Number       int           `db:"number" json:"number"`
	Seating      SeatingMethod `db:"seating" json:"seating"`
	Step         *int          `db:"seating_step" json:"seating_step"`
	AvoidRepeats bool          `db:"avoid_repeats" json:"avoid_repeats"`
	// Tables are the tour's tables imported from a file, participant
	// numbers east to north; nil without them. FileTables is how many.
	Tables     [][4]int `db:"tables" json:"-"`
	FileTables int      `db:"-" json:"file_tables"`
	// Seated: the tour's seating is published (or played); its plan no
	// longer matters.
	Seated bool `db:"-" json:"seated"`
}

// PlanTours is a seating for several tours: per tour number, its tables,
// each four participant numbers east to north.
type PlanTours map[int][][4]int

// SeatingFileInfo describes the tables imported from a file: their tours,
// the most tables in a tour, how many participant numbers they seat, and
// their repeat meetings among their own tours (as in TourState).
type SeatingFileInfo struct {
	Tours   []int `json:"tours"`
	Tables  int   `json:"tables"`
	Players int   `json:"players"`
	Repeats int   `json:"repeats"`
	MaxMet  int   `json:"max_met"`
}

// Participant is a seated-to-be participant: a numbered player who hasn't
// withdrawn.
type Participant struct {
	UserID int `db:"user_id"`
	Number int `db:"number"`
}

// SeatingPair is two players, the lower user ID first.
type SeatingPair [2]int

// SeatingStanding is a participant's totals from the finished tours.
type SeatingStanding struct {
	TP float64
	GP int
}

// SeatingInput is what a tour's seating is made from: the participants to
// seat (numbered, not withdrawn; a multiple of four), how many times each
// pair has already met in this tournament's tours, the standings so far
// (missing ones count as zero), the tour's number, its tables imported from
// a file (nil without them), and how many random tours the plan has after
// it (a random seating avoids repeats with them too).
type SeatingInput struct {
	Participants []Participant
	Meetings     map[SeatingPair]int
	Standings    map[int]SeatingStanding
	Tour         int
	Tables       [][4]int
	FutureRandom int
}

// TablesFunc makes a tour's tables: each table's four user IDs, east to
// north.
type TablesFunc func(input SeatingInput) ([][4]int, error)

// TourSeat is where a player sits in a tour, for their own table screen.
type TourSeat struct {
	Tour  int `json:"tour"`
	Table int `json:"table"`
}

// SeatChange is a referee seating someone at a seat of a tournament game:
// a substitute, or the seat's player back (SitterID == SeatPlayerID), from
// the log entry after SeenLogSize. Forfeit: the seat's player gets 0 TP and
// 0 GP for the game (only with a substitute).
type SeatChange struct {
	EventID      int
	GameID       int
	SeatPlayerID int
	SitterID     int
	Forfeit      bool
	SeenLogSize  int
}

// GameTimeLimit is a tournament game's time: when its time is up (the
// tour's start, plus its length, plus the time added to the table) and the
// tournament's rule for what happens then.
type GameTimeLimit struct {
	EndsAt     time.Time `db:"ends_at"`
	TimeUp     string    `db:"time_up"`
	ExtraHands int       `db:"extra_hands"`
}

// GameTimer is a tournament game's timer as a table shows it.
type GameTimer struct {
	EndsAt time.Time `json:"ends_at"`
	TimeUp string    `json:"time_up"`
	// LastHand is the last hand to play once the time is up (with
	// TimeUpFinishHand or TimeUpExtraHands); 0 before.
	LastHand int `json:"last_hand,omitempty"`
	// Stopped: the time is up with TimeUpStopNow; no more entries, the game
	// is to be finished.
	Stopped bool `json:"stopped,omitempty"`
}

// Limits of the time a referee adds to a table at once (the frontend
// mirrors them).
const (
	MinExtraMinutes = 1
	MaxExtraMinutes = 60
)

// LiveTour is a tournament's live tour as a player sees it: its ID and
// number, and whether the player's own game in it has ended (they may watch
// the other tables then).
type LiveTour struct {
	ID       int  `db:"id"`
	Number   int  `db:"number"`
	Finished bool `db:"finished"`
}
