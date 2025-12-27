package models

import "time"

const (
	DefaultDiffLingerDuration = 15 * time.Second
)

type GameHUD struct {
	Live  *GameLiveHUD   `json:"live,omitempty"`
	Final []PlayerResult `json:"final,omitempty"`
}

type GameLiveHUD struct {
	Round      Round       `json:"round"`
	PlayerHUDs []PlayerHUD `json:"player_huds"`
	Message    string      `json:"message,omitempty"`
}

type PlayerHUD struct {
	Name  string `json:"name"`
	Wind  string `json:"wind"`
	Score int    `json:"score"`
	Diff  string `json:"diff,omitempty"`
}
