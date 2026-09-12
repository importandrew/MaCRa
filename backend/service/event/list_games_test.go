package event

import (
	"context"
	"testing"

	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListGamesByStatus(t *testing.T) {
	const eventID = 7

	t.Run("reject aborted", func(t *testing.T) {
		svc := NewService(NewMockEventRepo(t), nil, nil)
		_, err := svc.ListGamesByStatus(context.Background(), eventID, models.GameAborted, 50, 0)
		assert.Error(t, err)
	})

	t.Run("reject created", func(t *testing.T) {
		svc := NewService(NewMockEventRepo(t), nil, nil)
		_, err := svc.ListGamesByStatus(context.Background(), eventID, models.GameCreated, 50, 0)
		assert.Error(t, err)
	})

	t.Run("empty live", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().ListGamesByStatus(mock.Anything, eventID, models.GameLive, 51, 0).
			Return(nil, nil)

		svc := NewService(repo, nil, nil)
		got, err := svc.ListGamesByStatus(context.Background(), eventID, models.GameLive, 50, 0)
		require.NoError(t, err)
		assert.Empty(t, got.Games)
		assert.False(t, got.HasMore)
	})

	t.Run("empty ended", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().ListGamesByStatus(mock.Anything, eventID, models.GameEnded, 51, 0).
			Return(nil, nil)

		svc := NewService(repo, nil, nil)
		got, err := svc.ListGamesByStatus(context.Background(), eventID, models.GameEnded, 50, 0)
		require.NoError(t, err)
		assert.Empty(t, got.Games)
		assert.False(t, got.HasMore)
	})
}
