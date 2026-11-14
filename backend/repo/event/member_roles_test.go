package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
	"github.com/importandrew/MaCRa/repo/database"
)

func TestSetMemberRoles(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "set roles", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestEvent(t, db) // owner 1: player and referee
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventPlayer, models.EventReferee}))

		require.NoError(t, repo.SetMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventReferee, models.EventObserver}))
		roles, err := repo.GetMemberRoles(ctx, eventID, 2)
		require.NoError(t, err)
		assert.ElementsMatch(t, []models.EventRole{models.EventReferee, models.EventObserver}, roles)

		require.NoError(t, repo.SetMemberRoles(ctx, eventID, 1, []models.EventRole{}))
		roles, err = repo.GetMemberRoles(ctx, eventID, 1)
		require.NoError(t, err)
		assert.Empty(t, roles, "the owner may be left without roles")

		roles, err = repo.GetMemberRoles(ctx, eventID, 2)
		require.NoError(t, err)
		assert.Len(t, roles, 2, "other members are untouched")
	})
}

// TestListMembers_ownerWithoutRoles checks that the owner stays in the member
// list after dropping every role.
func TestListMembers_ownerWithoutRoles(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	db.Test(t, "owner without roles", func(t *testing.T) {
		ctx := context.Background()
		repo := NewRepo(db)
		eventID := createTestEvent(t, db) // owner 1
		require.NoError(t, repo.AddMemberRoles(ctx, eventID, 2, []models.EventRole{models.EventPlayer}))
		require.NoError(t, repo.SetMemberRoles(ctx, eventID, 1, []models.EventRole{}))

		members, err := repo.ListMembers(ctx, eventID)
		require.NoError(t, err)
		assert.Equal(t, []models.EventMember{
			{UserID: 1, Name: "user1", Roles: []models.EventRole{}},
			{UserID: 2, Name: "user2", Roles: []models.EventRole{models.EventPlayer}},
		}, members)
	})
}
