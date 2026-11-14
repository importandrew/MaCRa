package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

func TestInvite(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "invite", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestEvent(t, db)
		createTestUser(t, db, 7)

		token, err := repo.InviteToken(ctx, eventID)
		require.NoError(t, err)
		assert.Nil(t, token, "off by default")

		first, second := "first-token", "second-token"
		require.NoError(t, repo.SetInviteToken(ctx, eventID, &first))
		require.NoError(t, repo.SetInviteToken(ctx, eventID, &second))
		_, err = repo.JoinByInvite(ctx, first, 7)
		assert.ErrorIs(t, err, models.ErrInviteNotFound, "a replaced link stops working")

		event, err := repo.EventByInviteToken(ctx, second)
		require.NoError(t, err)
		assert.Equal(t, eventID, event.ID)

		joined, err := repo.JoinByInvite(ctx, second, 7)
		require.NoError(t, err)
		assert.Equal(t, eventID, joined)
		roles, err := repo.GetMemberRoles(ctx, eventID, 7)
		require.NoError(t, err)
		assert.Equal(t, []models.EventRole{models.EventPlayer}, roles)
		_, err = repo.JoinByInvite(ctx, second, 7)
		assert.NoError(t, err, "joining twice is fine")

		require.NoError(t, repo.SetInviteToken(ctx, eventID, nil))
		_, err = repo.EventByInviteToken(ctx, second)
		assert.ErrorIs(t, err, models.ErrInviteNotFound, "turned off")

		require.NoError(t, repo.SetInviteToken(ctx, eventID, &first))
		require.NoError(t, repo.Finish(ctx, eventID))
		_, err = repo.JoinByInvite(ctx, first, 2)
		assert.ErrorIs(t, err, models.ErrEventAlreadyFinished)
		third := "third-token"
		assert.ErrorIs(t, repo.SetInviteToken(ctx, eventID, &third), models.ErrEventAlreadyFinished)
		assert.NoError(t, repo.SetInviteToken(ctx, eventID, nil), "turning it off is always fine")
	})
}
