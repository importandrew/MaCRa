package event

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
)

func TestSeatingFunc_unknown(t *testing.T) {
	_, err := seatingFunc(models.Seating{Method: "swiss-ish"})
	assert.ErrorIs(t, err, models.ErrInvalidSeating)
}

func TestSeatingFunc_intervalStep(t *testing.T) {
	input := models.SeatingInput{}
	for n := 1; n <= 8; n++ { // two tables
		input.Participants = append(input.Participants, models.Participant{UserID: 100 + n, Number: n})
	}
	step := func(k int) *int { return &k }

	_, err := seatingFunc(models.Seating{Method: models.SeatingRandom, Step: step(1)})
	assert.ErrorIs(t, err, models.ErrInvalidSeating, "a step for another method")

	for k, ok := range map[int]bool{0: false, 1: true, 2: true, 3: false} {
		makeTables, err := seatingFunc(models.Seating{Method: models.SeatingInterval, Step: step(k)})
		assert.NoError(t, err)
		_, err = makeTables(input)
		if ok {
			assert.NoError(t, err, "step %d", k)
		} else {
			assert.ErrorIs(t, err, models.ErrInvalidSeating, "step %d", k)
		}
	}
}
