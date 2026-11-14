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
	// MinRatingGames is how many ended games a player needs for a place in
	// the rating; players with fewer are listed below it without a place.
	// 1 lists everyone who has played.
	MinRatingGames int `json:"min_rating_games"`
}

// Limits of MinRatingGames (the frontend mirrors them).
const (
	MinMinRatingGames = 1
	MaxMinRatingGames = 100
)

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
// players too before the first round; everyone who has played is in the
// rating.
func DefaultEventSettings() EventSettings {
	return EventSettings{
		FinishGame:     FinishBySeated,
		AbortGame:      AbortBySeatedBeforeFirstRound,
		MinRatingGames: MinMinRatingGames,
	}
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

	if s.MinRatingGames < MinMinRatingGames || s.MinRatingGames > MaxMinRatingGames {
		return fmt.Errorf("%w: min_rating_games should be %d to %d, got %d",
			ErrInvalidEventSettings, MinMinRatingGames, MaxMinRatingGames, s.MinRatingGames)
	}

	return nil
}
