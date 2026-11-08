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
	// LogSize is the number of round log entries the outcome was computed
	// from; the repo refuses to save it if the log has changed since.
	LogSize int
}
