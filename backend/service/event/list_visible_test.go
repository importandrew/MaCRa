package event

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestListVisible(t *testing.T) {
	played := time.Date(2026, 11, 10, 18, 0, 0, 0, time.UTC)
	rows := []models.EventListRow{
		{
			ID: 1, OwnerUserID: 9, Visibility: models.EventVisibilityPrivate,
			Roles: []models.EventRole{models.EventPlayer}, LastGameAt: &played,
		},
		{ID: 2, OwnerUserID: 9, Visibility: models.EventVisibilityPrivate},
		{ID: 3, OwnerUserID: 9, Visibility: models.EventVisibilityPublic},
		{ID: 4, OwnerUserID: 5, Visibility: models.EventVisibilityPrivate},
	}

	t.Run("signed in", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().List(mock.Anything, 5).Return(rows, nil)

		got, err := NewService(repo, nil, nil).ListVisible(context.Background(), 5)

		require.NoError(t, err)
		// 2 is private without the caller; the order is the repo's.
		require.Len(t, got, 3)
		assert.Equal(t, []int{1, 3, 4}, []int{got[0].ID, got[1].ID, got[2].ID})
		assert.True(t, got[0].IsParticipating)
		assert.Equal(t, []models.EventRole{models.EventPlayer}, got[0].Roles)
		assert.Equal(t, &played, got[0].LastGameAt)
		assert.False(t, got[1].IsParticipating)
		assert.True(t, got[2].IsOwner, "the owner of a private event without roles sees it")
	})

	t.Run("visitor", func(t *testing.T) {
		repo := NewMockEventRepo(t)
		repo.EXPECT().List(mock.Anything, 0).Return(rows, nil)

		got, err := NewService(repo, nil, nil).ListVisible(context.Background(), 0)

		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, 3, got[0].ID)
		assert.False(t, got[0].IsParticipating)
		assert.Nil(t, got[0].Roles)
	})
}
