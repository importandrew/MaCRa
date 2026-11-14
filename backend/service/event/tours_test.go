package event

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
)

func TestSeatingFunc_unknown(t *testing.T) {
	_, err := seatingFunc("swiss-ish")
	assert.ErrorIs(t, err, models.ErrInvalidSeating)
}
