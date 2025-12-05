package models

import "fmt"

type Game struct{}

type GameNow struct {
	ID    int   `json:"id"`
	Round Round `json:"round"`

	Me    PlayerWindScore `json:"me"`
	Right PlayerWindScore `json:"right"`
	Opp   PlayerWindScore `json:"opp"`
	Left  PlayerWindScore `json:"left"`
}

type Round struct {
	Num   int    `json:"num"`
	Kanji string `json:"kanji"`
}

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

type NewGameResponse struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type GameResult struct {
	ID            int            `json:"id"`
	PlayerResults []PlayerResult `json:"player_results"`
}

type PlayerResult struct {
	Name string  `db:"name" json:"name"`
	ID   int     `db:"id" json:"id"`
	TP   float32 `db:"tp" json:"tp"`
	GP   int     `db:"gp" json:"gp"`
}
