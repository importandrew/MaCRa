package models

import (
	"fmt"
	"time"
)

type GameNow struct {
	ID     int        `json:"id"`
	Status GameStatus `json:"status"`
	Round  Round      `json:"round"`

	Me    PlayerWindScore `json:"me"`
	Right PlayerWindScore `json:"right"`
	Opp   PlayerWindScore `json:"opp"`
	Left  PlayerWindScore `json:"left"`
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
	ID               int            `json:"id"`
	FinishedAt       *time.Time     `json:"finished_at"`
	CanUndoLastRound bool           `json:"can_undo_last_round"`
	PlayerResults    []PlayerResult `json:"player_results"`
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

	PlayerScores []PlayerWindScore `json:"player_scores"`
}

type PlayerHistory struct {
	ID      int   `json:"id"`
	GameIDs []int `json:"game_ids"`
}

type EventGame struct {
	GameID  int `db:"id"`
	EventID int `db:"event_id"`
}

// Game is a row from the games table.
type Game struct {
	ID         int        `db:"id" json:"id"`
	Status     GameStatus `db:"status" json:"status"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	FinishedAt *time.Time `db:"finished_at" json:"finished_at"`
	EventID    int        `db:"event_id" json:"event_id"`
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
	ID      int                    `json:"id"`
	Status  GameStatus             `json:"status"`
	Round   *Round                 `json:"round"`
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
