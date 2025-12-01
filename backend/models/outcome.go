package models

type Outcome string

var (
	Deal Outcome = "deal-in"
	Self Outcome = "self"
	Draw Outcome = "draw"

	FalseMahjong Outcome = "false mahjong"
	Penalty      Outcome = "penalty"
)
