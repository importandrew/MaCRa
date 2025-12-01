package models

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
