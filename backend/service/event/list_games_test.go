package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestListGamesByStatus(t *testing.T) {
	const eventID = 7

	t.Run("reject aborted", func(t *testing.T) {
		svc := NewService(NewMockEventRepo(t), nil, nil)
		_, err := svc.ListGamesByStatus(context.Background(), eventID, 3, models.GameAborted, 50, 0)
		assert.Error(t, err)
	})

	t.Run("reject created", func(t *testing.T) {
		svc := NewService(NewMockEventRepo(t), nil, nil)
		_, err := svc.ListGamesByStatus(context.Background(), eventID, 3, models.GameCreated, 50, 0)
		assert.Error(t, err)
	})

	t.Run("empty live", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().Get(mock.Anything, eventID).Return(models.Event{ID: eventID, Type: models.EventRegular}, nil)
		repo.EXPECT().ListGamesByStatus(mock.Anything, eventID, models.GameLive, 0, 3, 51, 0).
			Return(nil, nil)

		svc := NewService(repo, nil, nil)
		got, err := svc.ListGamesByStatus(context.Background(), eventID, 3, models.GameLive, 50, 0)
		require.NoError(t, err)
		assert.Empty(t, got.Games)
		assert.False(t, got.HasMore)
	})

	t.Run("empty ended", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().ListGamesByStatus(mock.Anything, eventID, models.GameEnded, 0, 3, 51, 0).
			Return(nil, nil)

		svc := NewService(repo, nil, nil)
		got, err := svc.ListGamesByStatus(context.Background(), eventID, 3, models.GameEnded, 50, 0)
		require.NoError(t, err)
		assert.Empty(t, got.Games)
		assert.False(t, got.HasMore)
	})

	t.Run("live tour hidden from a player still playing", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().Get(mock.Anything, eventID).Return(models.Event{ID: eventID, Type: models.EventTournament}, nil)
		repo.EXPECT().GetMemberRoles(mock.Anything, eventID, 3).Return([]models.EventRole{models.EventPlayer}, nil)
		repo.EXPECT().LiveTour(mock.Anything, eventID, 3).Return(models.LiveTour{ID: 40, Number: 2}, nil)
		repo.EXPECT().ListGamesByStatus(mock.Anything, eventID, models.GameLive, 40, 3, 51, 0).
			Return(nil, nil)

		got, err := NewService(repo, nil, nil).ListGamesByStatus(context.Background(), eventID, 3, models.GameLive, 50, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, got.HiddenTour)
	})

	t.Run("live tour shown once the player's game ended", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().Get(mock.Anything, eventID).Return(models.Event{ID: eventID, Type: models.EventTournament}, nil)
		repo.EXPECT().GetMemberRoles(mock.Anything, eventID, 3).Return([]models.EventRole{models.EventPlayer}, nil)
		repo.EXPECT().LiveTour(mock.Anything, eventID, 3).Return(models.LiveTour{ID: 40, Number: 2, Finished: true}, nil)
		repo.EXPECT().ListGamesByStatus(mock.Anything, eventID, models.GameLive, 0, 3, 51, 0).
			Return(nil, nil)

		got, err := NewService(repo, nil, nil).ListGamesByStatus(context.Background(), eventID, 3, models.GameLive, 50, 0)
		require.NoError(t, err)
		assert.Zero(t, got.HiddenTour)
	})

	t.Run("staff see the live tour", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().Get(mock.Anything, eventID).Return(models.Event{ID: eventID, Type: models.EventTournament}, nil)
		repo.EXPECT().GetMemberRoles(mock.Anything, eventID, 3).Return([]models.EventRole{models.EventObserver}, nil)
		repo.EXPECT().ListGamesByStatus(mock.Anything, eventID, models.GameLive, 0, 3, 51, 0).
			Return(nil, nil)

		_, err := NewService(repo, nil, nil).ListGamesByStatus(context.Background(), eventID, 3, models.GameLive, 50, 0)
		require.NoError(t, err)
	})
}
