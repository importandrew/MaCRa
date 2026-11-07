package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerInAnotherGameError(t *testing.T) {
	err := fmt.Errorf("undo: %w", PlayerInAnotherGameError{Name: "Ivan"})

	assert.ErrorIs(t, err, ErrUndoFinalRound)
	assert.NotErrorIs(t, err, ErrLogChanged)
	assert.EqualError(t, PlayerInAnotherGameError{Name: "Ivan"}, "Ivan has already started another game")

	var target PlayerInAnotherGameError
	assert.ErrorAs(t, err, &target)
	assert.Equal(t, "Ivan", target.Name)
}
