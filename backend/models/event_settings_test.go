package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventSettingsValidate(t *testing.T) {
	testData := map[string]struct {
		settings EventSettings
		valid    bool
	}{
		"default":                  {DefaultEventSettings(), true},
		"referees only":            {EventSettings{FinishGame: FinishByReferees, AbortGame: AbortByReferees, UndoRound: UndoByReferees, MinRatingGames: 1}, true},
		"seated abort at any time": {EventSettings{FinishGame: FinishBySeated, AbortGame: AbortBySeated, UndoRound: UndoByReferees, MinRatingGames: 100}, true},
		"unknown finish":           {EventSettings{FinishGame: "anyone", AbortGame: AbortByReferees, UndoRound: UndoByReferees, MinRatingGames: 1}, false},
		"unknown abort":            {EventSettings{FinishGame: FinishBySeated, AbortGame: "anyone", UndoRound: UndoByReferees, MinRatingGames: 1}, false},
		"min rating games 0":       {EventSettings{FinishGame: FinishBySeated, AbortGame: AbortBySeated, UndoRound: UndoByReferees, MinRatingGames: 0}, false},
		"min rating games 101":     {EventSettings{FinishGame: FinishBySeated, AbortGame: AbortBySeated, UndoRound: UndoByReferees, MinRatingGames: 101}, false},
		"undo by seated":           {EventSettings{FinishGame: FinishBySeated, AbortGame: AbortBySeated, UndoRound: UndoBySeated, MinRatingGames: 1}, true},
		"unknown undo":             {EventSettings{FinishGame: FinishBySeated, AbortGame: AbortBySeated, UndoRound: "anyone", MinRatingGames: 1}, false},
		"empty":                    {EventSettings{}, false},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			err := testCase.settings.Validate()
			if testCase.valid {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, ErrInvalidEventSettings)
			}
		})
	}
}
