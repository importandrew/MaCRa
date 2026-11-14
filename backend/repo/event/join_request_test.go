package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

func TestJoinRequests(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "join requests", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestEvent(t, db) // public, players 1–4
		for id := 5; id <= 9; id++ {
			createTestUser(t, db, id)
		}

		event, err := repo.Get(ctx, eventID)
		require.NoError(t, err)
		assert.False(t, event.AcceptRequests, "off by default")
		assert.Nil(t, event.PlayerLimit, "no limit by default")
		assert.ErrorIs(t, repo.CreateJoinRequest(ctx, eventID, 5), models.ErrRequestsClosed, "switched off")

		limit := 6
		require.NoError(t, repo.UpdateJoining(ctx, eventID, models.JoiningSettings{AcceptRequests: true, PlayerLimit: &limit}))
		event, err = repo.Get(ctx, eventID)
		require.NoError(t, err)
		assert.True(t, event.AcceptRequests)
		assert.Equal(t, &limit, event.PlayerLimit)

		assert.ErrorIs(t, repo.CreateJoinRequest(ctx, eventID, 2), models.ErrAlreadyMember)
		for id := 5; id <= 8; id++ {
			require.NoError(t, repo.CreateJoinRequest(ctx, eventID, id))
		}
		require.NoError(t, repo.CreateJoinRequest(ctx, eventID, 5), "asking again while pending changes nothing")

		requests, err := repo.ListJoinRequests(ctx, eventID)
		require.NoError(t, err)
		require.Len(t, requests, 4)
		assert.Equal(t, 5, requests[0].UserID, "oldest first")
		assert.Equal(t, "user5", requests[0].Name)

		require.NoError(t, repo.WithdrawJoinRequest(ctx, eventID, 8))
		assert.ErrorIs(t, repo.WithdrawJoinRequest(ctx, eventID, 8), models.ErrRequestNotFound)

		require.NoError(t, repo.DeclineJoinRequest(ctx, eventID, 7))
		assert.ErrorIs(t, repo.CreateJoinRequest(ctx, eventID, 7), models.ErrRequestDeclined, "can't ask again")
		assert.ErrorIs(t, repo.AcceptJoinRequest(ctx, eventID, 7), models.ErrRequestNotFound, "a declined request can't be accepted")
		status, err := repo.JoinRequestStatus(ctx, eventID, 7)
		require.NoError(t, err)
		assert.Equal(t, models.JoinRequestDeclined, status)

		// Added by hand: the pending request goes.
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 6, []models.EventRole{models.EventObserver}))
		status, err = repo.JoinRequestStatus(ctx, eventID, 6)
		require.NoError(t, err)
		assert.Empty(t, status)

		require.NoError(t, repo.AcceptJoinRequest(ctx, eventID, 5))
		roles, err := repo.GetMemberRoles(ctx, eventID, 5)
		require.NoError(t, err)
		assert.Equal(t, []models.EventRole{models.EventPlayer}, roles)
		count, err := repo.CountJoinRequests(ctx, eventID)
		require.NoError(t, err)
		assert.Zero(t, count)

		// Five players of six: one place left, the invite link takes it.
		token := "invite-token"
		require.NoError(t, repo.SetInviteToken(ctx, eventID, &token))
		require.NoError(t, repo.CreateJoinRequest(ctx, eventID, 9))
		_, err = repo.JoinByInvite(ctx, token, 8)
		require.NoError(t, err)
		assert.ErrorIs(t, repo.AcceptJoinRequest(ctx, eventID, 9), models.ErrEventFull)
		_, err = repo.JoinByInvite(ctx, token, 9)
		assert.ErrorIs(t, err, models.ErrEventFull)
		_, err = repo.JoinByInvite(ctx, token, 5)
		assert.NoError(t, err, "a player opening the link again is fine when full")
		require.NoError(t, repo.WithdrawJoinRequest(ctx, eventID, 9))
		assert.ErrorIs(t, repo.CreateJoinRequest(ctx, eventID, 9), models.ErrEventFull)
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 9, []models.EventRole{models.EventPlayer}), "adding by hand isn't limited")

		require.NoError(t, repo.UpdateJoining(ctx, eventID, models.JoiningSettings{AcceptRequests: true}))
		require.NoError(t, repo.Finish(ctx, eventID))
		assert.ErrorIs(t, repo.CreateJoinRequest(ctx, eventID, 7), models.ErrRequestsClosed)
		assert.ErrorIs(t, repo.UpdateJoining(ctx, eventID, models.JoiningSettings{}), models.ErrEventAlreadyFinished)
	})
}

func TestCreateJoinRequest_privateEvent(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "private", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestEvent(t, db)
		createTestUser(t, db, 5)
		require.NoError(t, repo.UpdateJoining(ctx, eventID, models.JoiningSettings{AcceptRequests: true}))
		require.NoError(t, repo.CreateJoinRequest(ctx, eventID, 5))

		event, err := repo.Get(ctx, eventID)
		require.NoError(t, err)
		event.Visibility = models.EventVisibilityPrivate
		require.NoError(t, repo.Update(ctx, event))

		createTestUser(t, db, 6)
		assert.ErrorIs(t, repo.CreateJoinRequest(ctx, eventID, 6), models.ErrRequestsClosed)
		assert.NoError(t, repo.AcceptJoinRequest(ctx, eventID, 5), "the owner still answers pending requests")
	})
}
