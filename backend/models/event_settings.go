package models

import (
	"fmt"
	"slices"
)

// EventSettings are an event's settings that don't affect the score: who may
// do what with a game. Unlike GameRules they can change at any time.
type EventSettings struct {
	// FinishGame is who may finish a live game early (it counts):
	// FinishBySeated or FinishByReferees.
	FinishGame string `json:"finish_game"`
	// AbortGame is who may abort a live game (it doesn't count):
	// AbortByReferees, AbortBySeatedBeforeFirstRound or AbortBySeated.
	AbortGame string `json:"abort_game"`
}

// Who may finish a game early. Referees always may.
const (
	// FinishBySeated: the seated players and referees.
	FinishBySeated = "seated"
	// FinishByReferees: referees only.
	FinishByReferees = "referees"
)

// Who may abort a game. Referees always may.
const (
	// AbortByReferees: referees only.
	AbortByReferees = "referees"
	// AbortBySeatedBeforeFirstRound: also the seated players, while no round
	// has been entered (a game started by mistake).
	AbortBySeatedBeforeFirstRound = "seated_before_first_round"
	// AbortBySeated: also the seated players, at any time.
	AbortBySeated = "seated"
)

var (
	ValidFinishGame = []string{FinishBySeated, FinishByReferees}
	ValidAbortGame  = []string{AbortByReferees, AbortBySeatedBeforeFirstRound, AbortBySeated}
)

// DefaultEventSettings returns the settings events start with: seated
// players and referees finish games, referees abort them, and the seated
// players too before the first round.
func DefaultEventSettings() EventSettings {
	return EventSettings{FinishGame: FinishBySeated, AbortGame: AbortBySeatedBeforeFirstRound}
}

// Validate returns an error wrapping ErrInvalidEventSettings if a value is
// not one of its options.
func (s EventSettings) Validate() error {
	if !slices.Contains(ValidFinishGame, s.FinishGame) {
		return fmt.Errorf("%w: finish_game should be one of %v, got %q",
			ErrInvalidEventSettings, ValidFinishGame, s.FinishGame)
	}

	if !slices.Contains(ValidAbortGame, s.AbortGame) {
		return fmt.Errorf("%w: abort_game should be one of %v, got %q",
			ErrInvalidEventSettings, ValidAbortGame, s.AbortGame)
	}

	return nil
}
