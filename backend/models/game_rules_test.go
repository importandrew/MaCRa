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
		"standard":              {StandardGameRules(), true},
		"8 hands":               {GameRules{Hands: 8, MinHand: 8}, true},
		"4 hands":               {GameRules{Hands: 4, MinHand: 8}, true},
		"1 hand":                {GameRules{Hands: 1, MinHand: 8}, true},
		"highest minimum hand":  {GameRules{Hands: 16, MinHand: 88}, true},
		"12 hands":              {GameRules{Hands: 12, MinHand: 8}, false},
		"minimum hand below 8":  {GameRules{Hands: 16, MinHand: 7}, false},
		"minimum hand above 88": {GameRules{Hands: 16, MinHand: 89}, false},
		"empty":                 {GameRules{}, false},
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
