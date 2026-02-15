package models

import "time"

const UndoFinalRoundTimeout = 1 * time.Hour

type FinalGameData struct {
	GameID int
	East   FinalPlayerData
	South  FinalPlayerData
	West   FinalPlayerData
	North  FinalPlayerData
}

type FinalPlayerData struct {
	ID int
	GP int
	TP float32
}
