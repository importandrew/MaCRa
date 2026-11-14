package models

import (
	"fmt"
	"slices"
)

// GameRules are the rules that affect a game's score. An event has them,
// frozen after its first game, and each game keeps a copy of the rules it
// started with: it is always replayed by that copy.
type GameRules struct {
	// Hands is the number of hands: 16 (East, South, West and North rounds),
	// 8, 4 or 1.
	Hands int `json:"hands"`
	// MinHand is the lowest hand value a player may win with.
	MinHand int `json:"min_hand"`
}

// Limits of the game rules (the frontend mirrors them).
const (
	MinMinHand = 8
	MaxMinHand = 88
)

// ValidHands is the allowed numbers of hands in a game.
var ValidHands = []int{1, 4, 8, 16}

// FalseMahjongAmounts is what a false mahjong may cost; the other three
// players get a third each.
var FalseMahjongAmounts = []int{30, 60}

// StandardGameRules returns the MCR standard rules: 16 hands, minimum hand 8.
// Events start with them.
//
//nolint:mnd
func StandardGameRules() GameRules {
	return GameRules{Hands: 16, MinHand: 8}
}

// Validate returns an error wrapping ErrInvalidGameRules if a value is out
// of its limits.
func (r GameRules) Validate() error {
	if !slices.Contains(ValidHands, r.Hands) {
		return fmt.Errorf("%w: hands should be one of %v, got %d", ErrInvalidGameRules, ValidHands, r.Hands)
	}

	if r.MinHand < MinMinHand || r.MinHand > MaxMinHand {
		return fmt.Errorf("%w: minimum hand should be %d to %d, got %d",
			ErrInvalidGameRules, MinMinHand, MaxMinHand, r.MinHand)
	}

	return nil
}
