package models

type GameStatus string

const (
	GameCreated GameStatus = "created"
	GameLive    GameStatus = "live"
	GameEnded   GameStatus = "ended"
	GameAborted GameStatus = "aborted"
)
