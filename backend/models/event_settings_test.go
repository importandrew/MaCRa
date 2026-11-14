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
		"referees only":            {EventSettings{FinishGame: FinishByReferees, AbortGame: AbortByReferees}, true},
		"seated abort at any time": {EventSettings{FinishGame: FinishBySeated, AbortGame: AbortBySeated}, true},
		"unknown finish":           {EventSettings{FinishGame: "anyone", AbortGame: AbortByReferees}, false},
		"unknown abort":            {EventSettings{FinishGame: FinishBySeated, AbortGame: "anyone"}, false},
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
