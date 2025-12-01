package models

type RoundData struct {
	Round           int
	Outcome         Outcome
	WinnerID        int
	LoserID         int
	Value           int
	EastStartScore  int
	SouthStartScore int
	WestStartScore  int
	NorthStartScore int
	ForceNextRound  bool
}
