package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameRulesValidate(t *testing.T) {
	testData := map[string]struct {
		rules GameRules
		valid bool
	}{
		"standard":                 {StandardGameRules(), true},
		"short":                    {GameRules{Hands: 8, MinHand: 24, FalseMahjong: [2]int{30, 60}}, true},
		"limits":                   {GameRules{Hands: 16, MinHand: 88, FalseMahjong: [2]int{3, 300}}, true},
		"12 hands":                 {GameRules{Hands: 12, MinHand: 8, FalseMahjong: [2]int{30, 60}}, false},
		"minimum hand below 8":     {GameRules{Hands: 16, MinHand: 7, FalseMahjong: [2]int{30, 60}}, false},
		"minimum hand above 88":    {GameRules{Hands: 16, MinHand: 89, FalseMahjong: [2]int{30, 60}}, false},
		"false mahjong zero":       {GameRules{Hands: 16, MinHand: 8, FalseMahjong: [2]int{0, 60}}, false},
		"false mahjong not by 3":   {GameRules{Hands: 16, MinHand: 8, FalseMahjong: [2]int{30, 61}}, false},
		"false mahjong above 300":  {GameRules{Hands: 16, MinHand: 8, FalseMahjong: [2]int{30, 303}}, false},
		"false mahjong equal":      {GameRules{Hands: 16, MinHand: 8, FalseMahjong: [2]int{30, 30}}, false},
		"false mahjong going down": {GameRules{Hands: 16, MinHand: 8, FalseMahjong: [2]int{60, 30}}, false},
		"empty":                    {GameRules{}, false},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			err := testCase.rules.Validate()
			if testCase.valid {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, ErrInvalidGameRules)
			}
		})
	}
}
