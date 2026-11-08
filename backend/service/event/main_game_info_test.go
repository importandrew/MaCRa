package event

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
)

func TestGameFound(t *testing.T) {
	dbErr := errors.New("database error")

	testData := map[string]struct {
		err           error
		expectedFound bool
		expectedErr   error
	}{
		"found":              {err: nil, expectedFound: true},
		"not found":          {err: models.ErrGameNotFound},
		"wrapped not found":  {err: fmt.Errorf("get current game: %w", models.ErrGameNotFound)},
		"other error passed": {err: dbErr, expectedErr: dbErr},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			found, err := gameFound(testCase.err)

			assert.Equal(t, testCase.expectedFound, found)
			assert.ErrorIs(t, err, testCase.expectedErr)
		})
	}
}
