package models

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithData(t *testing.T) {
	err := fmt.Errorf("undo: %w", WithData(ErrPlayerInAnotherGame, ErrorData{"name": "Ivan"}))

	assert.ErrorIs(t, err, ErrPlayerInAnotherGame)
	assert.NotErrorIs(t, err, ErrLogChanged)
	assert.EqualError(t, err, "undo: player is in another live game map[name:Ivan]")

	dataErr, ok := errors.AsType[*DataError](err)
	require.True(t, ok)
	assert.Equal(t, ErrorData{"name": "Ivan"}, dataErr.Data)
}

// TestErrorsCodesUnique checks that no two errors share a code: the frontend
// tells them apart by it.
func TestErrorsCodesUnique(t *testing.T) {
	seen := map[string]error{ErrInternal.Code: nil}
	for err, apiErr := range Errors {
		require.NotEmpty(t, apiErr.Code, err)
		prev, ok := seen[apiErr.Code]
		assert.False(t, ok, "%q: %v and %v", apiErr.Code, prev, err)
		seen[apiErr.Code] = err
	}
}
