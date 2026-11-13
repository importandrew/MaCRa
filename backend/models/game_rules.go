package models

import (
	"errors"
	"fmt"
	"slices"
)

// GameRules are the rules that affect a game's score. An event has them,
// frozen after its first game, and each game keeps a copy of the rules it
// started with: it is always replayed by that copy.
type GameRules struct {
	// Hands is the number of hands: 16 (East, South, West and North rounds)
	// or 8 (East and South).
	Hands int `json:"hands"`
	// MinHand is the lowest hand value a player may win with.
	MinHand int `json:"min_hand"`
	// FalseMahjong is the two amounts a false mahjong may cost, the smaller
	// first; the other three players get a third each.
	FalseMahjong [2]int `json:"false_mahjong"`
}

// Limits of the game rules (the frontend mirrors them).
const (
	MinMinHand          = 8
	MaxMinHand          = 88
	MaxFalseMahjong     = 300
	falseMahjongPlayers = 3
)

// ValidHands is the allowed numbers of hands in a game.
var ValidHands = []int{8, 16}

// StandardGameRules returns the MCR standard rules: 16 hands, minimum hand 8,
// false mahjong 30 or 60. Events start with them.
//
//nolint:mnd
func StandardGameRules() GameRules {
	return GameRules{Hands: 16, MinHand: 8, FalseMahjong: [2]int{30, 60}}
}

// Validate returns an error wrapping ErrInvalidGameRules if a value is out
// of its limits.
func (r GameRules) Validate() error {
	if err := r.validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidGameRules, err)
	}
	return nil
}

func (r GameRules) validate() error {
	if !slices.Contains(ValidHands, r.Hands) {
		return fmt.Errorf("hands should be one of %v, got %d", ValidHands, r.Hands)
	}

	if r.MinHand < MinMinHand || r.MinHand > MaxMinHand {
		return fmt.Errorf("minimum hand should be %d to %d, got %d", MinMinHand, MaxMinHand, r.MinHand)
	}

	for _, value := range r.FalseMahjong {
		if value <= 0 || value > MaxFalseMahjong || value%falseMahjongPlayers != 0 {
			return fmt.Errorf("false mahjong should be a multiple of %d up to %d, got %d",
				falseMahjongPlayers, MaxFalseMahjong, value)
		}
	}

	if r.FalseMahjong[0] >= r.FalseMahjong[1] {
		return errors.New("false mahjong amounts should go up")
	}

	return nil
}
