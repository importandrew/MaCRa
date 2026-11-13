package models

import (
	"fmt"
	"time"
)

type GameNow struct {
	ID     int        `json:"id"`
	Status GameStatus `json:"status"`
	Round  Round      `json:"round"`
	// LogSize is the number of round log entries; a new entry sends it back.
	LogSize int `json:"log_size"`
	// Rules are the rules the game is played by.
	Rules GameRules `json:"rules"`

	Me    PlayerWindScore `json:"me"`
	Right PlayerWindScore `json:"right"`
	Opp   PlayerWindScore `json:"opp"`
	Left  PlayerWindScore `json:"left"`

	// Log is the round log of the same state as the table, so the table
	// screen and its log can't disagree.
	Log GameLog `json:"log"`
}

type Round struct {
	Num   int    `json:"num"`
	Kanji string `json:"kanji"`
}

//nolint:mnd
func RoundKanji(round int) string {
	wind := Wind((round - 1) / 4).Kanji()
	num := (round-1)%4 + 1

	return fmt.Sprintf("%s%d", wind, num)
}

type PlayerWindScore struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Wind  string `json:"wind"`
	Score int    `json:"score"`
}

type GameResult struct {
	ID            int            `json:"id"`
	FinishedAt    *time.Time     `json:"finished_at"`
	PlayerResults []PlayerResult `json:"player_results"`
}

type PlayerResult struct {
	Name string  `db:"name" json:"name"`
	ID   int     `db:"id" json:"id"`
	TP   float32 `db:"tp" json:"tp"`
	GP   int     `db:"gp" json:"gp"`
}

type GameView struct {
	ID     int        `json:"id"`
	Status GameStatus `json:"status"`
	Round  Round      `json:"round"`
	// LogSize is the number of round log entries, as in GameNow.
	LogSize int       `json:"log_size"`
	Rules   GameRules `json:"rules"`

	PlayerScores []PlayerWindScore `json:"player_scores"`

	// Log is the round log of the same state, as in GameNow.
	Log GameLog `json:"log"`
}

type MainGameInfo struct {
	HasLast    bool `json:"has_last"`
	HasCurrent bool `json:"has_current"`
}

// EventGames is the paginated response for event games listing.
type EventGames struct {
	Games   []EventGameItem `json:"games"`
	HasMore bool            `json:"has_more"`
}

// EventGameItem is one game row in GET /api/event/:eventId/games.
type EventGameItem struct {
	ID     int        `json:"id"`
	Status GameStatus `json:"status"`
	Round  *Round     `json:"round"`
	// Hands is the number of hands in a live game, for its progress.
	Hands   int                    `json:"hands,omitempty"`
	Players []EventGamePlayerScore `json:"players"`
}

// EventGamePlayerScore is a player preview on an event games list item.
type EventGamePlayerScore struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Wind string `json:"wind,omitempty"`
	GP   int    `json:"gp"`
}

// EventGameItemRow is a DB row used while building EventGameListItem.
type EventGameItemRow struct {
	ID     int        `db:"id"`
	Status GameStatus `db:"status"`
}

// GameSummary is what the capability rules need to know about a game.
type GameSummary struct {
	Status     GameStatus `db:"status"`
	FinishedAt *time.Time `db:"finished_at"`
	LogSize    int        `db:"log_size"`
}

// GameMe is the caller's capabilities for one game
// (GET /event/:id/game/:gameId/me).
type GameMe struct {
	Capabilities []Capability `json:"capabilities"`
}
