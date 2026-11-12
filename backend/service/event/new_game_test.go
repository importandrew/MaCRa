package event

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestNewGame_missingPlayerRoleNamesThePlayer(t *testing.T) {
	repo := NewMockEventRepo(t)
	players := NewMockPlayerRepo(t)

	repo.EXPECT().Get(mock.Anything, 1).Return(models.Event{ID: 1}, nil)
	repo.EXPECT().GetMemberRoles(mock.Anything, 1, 10).Return([]models.EventRole{models.EventPlayer}, nil)
	repo.EXPECT().GetMemberRoles(mock.Anything, 1, 11).Return([]models.EventRole{models.EventReferee}, nil)
	players.EXPECT().Get(mock.Anything, 11).Return(models.Player{ID: 11, Name: "Ivan"}, nil)

	_, err := NewService(repo, players, nil).NewGame(context.Background(), 1, []int{10, 11, 12, 13})

	require.ErrorIs(t, err, models.ErrMissingPlayerRole)
	dataErr, ok := errors.AsType[*models.DataError](err)
	require.True(t, ok)
	assert.Equal(t, models.ErrorData{"name": "Ivan"}, dataErr.Data)
	assert.Contains(t, err.Error(), "player 11", "the log keeps the ID")
}
